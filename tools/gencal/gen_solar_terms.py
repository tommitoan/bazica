"""Generate the 24 solar terms per year from the Sun's apparent position.

A solar term is the instant the Sun's apparent ecliptic longitude of date
reaches a multiple of 15 degrees. Positions come from the JPL DE440 ephemeris
through skyfield, which also supplies the TT-UT1 difference (delta T).

The output has the same shape as bazica's data/solar-term.json so the library
can embed it unchanged.
"""

from __future__ import annotations

import argparse
import json
import sys
from dataclasses import dataclass
from datetime import datetime, timedelta, timezone
from pathlib import Path

import numpy as np
from skyfield import almanac, api

# Index = longitude / 15 degrees. The names match the keys of solar-term.json.
TERM_NAMES: tuple[str, ...] = (
    "spring_equinox",  # 0
    "pure_brightness",  # 15
    "grain_rain",  # 30
    "start_of_summer",  # 45
    "grain_buds",  # 60
    "grain_in_ear",  # 75
    "summer_solstice",  # 90
    "minor_heat",  # 105
    "major_heat",  # 120
    "start_of_autumn",  # 135
    "end_of_heat",  # 150
    "white_dew",  # 165
    "autumn_equinox",  # 180
    "cold_dew",  # 195
    "frost",  # 210
    "start_of_winter",  # 225
    "minor_snow",  # 240
    "major_snow",  # 255
    "winter_solstice",  # 270
    "minor_cold",  # 285
    "major_cold",  # 300
    "start_of_spring",  # 315
    "spring_showers",  # 330
    "awakening_of_insects",  # 345
)

# The fastest the Sun moves is about 1.02 degrees a day, so two 15-degree
# crossings are never closer than about 14.7 days; a 10-day step cannot skip one.
_SEARCH_STEP_DAYS = 10.0
# find_discrete stops bisecting when the bracket is narrower than this (days).
# A Julian date held as a float64 resolves about 40 microseconds, so asking for
# less never terminates and exhausts memory; 0.1 ms is below the millisecond we
# keep and above that floor.
_EPSILON_DAYS = 1e-4 / 86400.0

# Crossings closer than this are one crossing reported more than once.
_MIN_SEPARATION_SECONDS = 1.0

# A found crossing must sit this close to a multiple of 15 degrees (about 9 seconds
# of the Sun's motion); anything else means the search returned a wrong instant.
_MAX_RESIDUAL_DEGREES = 1e-4

_TIMESTAMP_FORMAT = "%Y-%m-%d %H:%M:%S"


@dataclass(frozen=True)
class SolarTerm:
    """One solar term instant in UTC, rounded to the millisecond."""

    name: str
    instant: datetime


def format_timestamp(instant: datetime) -> str:
    """Format an instant the way solar-term.json does: millisecond precision, +00:00."""
    millis = instant.microsecond // 1000
    return f"{instant.strftime(_TIMESTAMP_FORMAT)}.{millis:03d}+00:00"


def parse_timestamp(text: str) -> datetime:
    """Parse a timestamp from solar-term.json into an aware UTC datetime."""
    return datetime.strptime(text, "%Y-%m-%d %H:%M:%S.%f%z").astimezone(timezone.utc)


def _round_to_millisecond(instant: datetime) -> datetime:
    rounded = instant + timedelta(microseconds=500)
    return rounded.replace(microsecond=(rounded.microsecond // 1000) * 1000)


def _term_index(longitude_degrees: np.ndarray) -> np.ndarray:
    return np.floor(longitude_degrees / 15.0).astype(int) % 24


def find_solar_terms(
    ephemeris_path: Path, first_year: int, last_year: int
) -> list[SolarTerm]:
    """Return every solar term whose UTC instant falls in first_year..last_year.

    Terms come back in chronological order. Each year in the range is complete:
    the search window is the whole of first_year through last_year in UTC.
    """
    if first_year > last_year:
        raise ValueError(f"first_year {first_year} is after last_year {last_year}")

    ts = api.load.timescale()
    eph = api.load_file(str(ephemeris_path))
    earth, sun = eph["earth"], eph["sun"]

    def term_at(t) -> np.ndarray:
        _, lon, _ = earth.at(t).observe(sun).apparent().ecliptic_latlon(epoch="date")
        return _term_index(lon.degrees)

    term_at.step_days = _SEARCH_STEP_DAYS

    # Start slightly before the window so a term at the very start is a
    # transition and not the initial state.
    t0 = ts.utc(first_year - 1, 12, 1)
    t1 = ts.utc(last_year + 1, 1, 1)

    times, _ = almanac.find_discrete(t0, t1, term_at, epsilon=_EPSILON_DAYS)

    # At float resolution the step function can flicker across a crossing and
    # report it several times within microseconds; real terms are two weeks apart.
    seconds = times.tt * 86400.0
    keep = np.concatenate(([True], np.diff(seconds) > _MIN_SEPARATION_SECONDS))
    times = times[keep]

    # find_discrete reports the function value at the bracket's end, which can
    # still be the value from before the crossing when the bracket is as narrow
    # as the float resolution of a Julian date. Name each crossing from the
    # longitude at the instant itself, to the nearest multiple of 15 degrees.
    _, lon, _ = earth.at(times).observe(sun).apparent().ecliptic_latlon(epoch="date")
    steps = lon.degrees / 15.0
    indexes = np.rint(steps).astype(int) % 24
    residual_degrees = np.abs(steps - np.rint(steps)) * 15.0
    if residual_degrees.size and residual_degrees.max() > _MAX_RESIDUAL_DEGREES:
        raise ValueError(
            f"a crossing is {residual_degrees.max():.6f} degrees from a multiple of 15"
        )

    lower = datetime(first_year, 1, 1, tzinfo=timezone.utc)
    upper = datetime(last_year + 1, 1, 1, tzinfo=timezone.utc)
    terms: list[SolarTerm] = []
    for t, index in zip(times, indexes):
        instant = _round_to_millisecond(t.utc_datetime())
        if lower <= instant < upper:
            terms.append(SolarTerm(TERM_NAMES[int(index)], instant))
    return terms


def terms_by_year(terms: list[SolarTerm]) -> dict[int, dict[str, datetime]]:
    """Group terms by the UTC calendar year of their instant."""
    grouped: dict[int, dict[str, datetime]] = {}
    for term in terms:
        year_terms = grouped.setdefault(term.instant.year, {})
        if term.name in year_terms:
            raise ValueError(f"{term.name} occurs twice in {term.instant.year}")
        year_terms[term.name] = term.instant
    for year, year_terms in grouped.items():
        if len(year_terms) != len(TERM_NAMES):
            raise ValueError(f"{year} has {len(year_terms)} terms, want {len(TERM_NAMES)}")
    return grouped


def to_json_document(grouped: dict[int, dict[str, datetime]]) -> dict[str, dict]:
    """Build the solar-term.json structure (years ascending, term names sorted)."""
    return {
        str(year): {
            "year": str(year),
            "data": {name: format_timestamp(year_terms[name]) for name in sorted(year_terms)},
        }
        for year, year_terms in sorted(grouped.items())
    }


def write_json(document: dict, out_path: Path) -> None:
    """Write the document byte-for-byte in the layout of the bundled file (no trailing newline)."""
    out_path.write_text(json.dumps(document, indent=2, ensure_ascii=False), encoding="utf-8")


@dataclass(frozen=True)
class Difference:
    """A term where the generated instant differs from the bundled one."""

    year: int
    name: str
    seconds: float


def compare_with_bundled(
    generated: dict[int, dict[str, datetime]], bundled_document: dict
) -> tuple[list[Difference], list[int]]:
    """Return every per-term difference over the shared years, and those years."""
    shared_years = sorted(y for y in generated if str(y) in bundled_document)
    differences: list[Difference] = []
    for year in shared_years:
        bundled_terms = bundled_document[str(year)]["data"]
        for name, instant in generated[year].items():
            bundled_instant = parse_timestamp(bundled_terms[name])
            delta = (instant - bundled_instant).total_seconds()
            differences.append(Difference(year, name, delta))
    return differences, shared_years


def summarise(differences: list[Difference]) -> dict[str, float]:
    """Mean and maximum absolute difference, in seconds."""
    magnitudes = np.array([abs(d.seconds) for d in differences])
    return {
        "terms": float(len(magnitudes)),
        "mean_abs_s": float(magnitudes.mean()),
        "median_abs_s": float(np.median(magnitudes)),
        "max_abs_s": float(magnitudes.max()),
    }


def _cmd_generate(args: argparse.Namespace) -> int:
    terms = find_solar_terms(args.ephemeris, args.first_year, args.last_year)
    document = to_json_document(terms_by_year(terms))
    write_json(document, args.out)
    print(f"wrote {len(document)} years to {args.out}")
    return 0


def _cmd_validate(args: argparse.Namespace) -> int:
    bundled = json.loads(args.bundled.read_text(encoding="utf-8"))
    bundled_years = sorted(int(y) for y in bundled)
    terms = find_solar_terms(args.ephemeris, bundled_years[0], bundled_years[-1])
    differences, years = compare_with_bundled(terms_by_year(terms), bundled)
    stats = summarise(differences)
    print(f"years compared: {years[0]}-{years[-1]} ({len(years)} years, {int(stats['terms'])} terms)")
    print(
        f"mean |diff| {stats['mean_abs_s']:.3f} s, median {stats['median_abs_s']:.3f} s, "
        f"max {stats['max_abs_s']:.3f} s"
    )
    worst = sorted(differences, key=lambda d: -abs(d.seconds))[: args.top]
    print(f"largest {len(worst)} differences (generated - bundled):")
    for d in worst:
        print(f"  {d.year} {d.name}: {d.seconds:+.3f} s")
    over = [d for d in differences if abs(d.seconds) > args.limit]
    print(f"terms differing by more than {args.limit:g} s: {len(over)}")
    return 1 if over else 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="command", required=True)

    gen = sub.add_parser("generate", help="write solar-term.json for a range of years")
    gen.add_argument("--ephemeris", type=Path, required=True, help="path to de440.bsp")
    gen.add_argument("--first-year", type=int, required=True)
    gen.add_argument("--last-year", type=int, required=True)
    gen.add_argument("--out", type=Path, required=True)
    gen.set_defaults(func=_cmd_generate)

    val = sub.add_parser("validate", help="compare the generator with a bundled solar-term.json")
    val.add_argument("--ephemeris", type=Path, required=True, help="path to de440.bsp")
    val.add_argument("--bundled", type=Path, required=True, help="path to data/solar-term.json")
    val.add_argument("--top", type=int, default=10, help="how many worst terms to list")
    val.add_argument("--limit", type=float, default=60.0, help="seconds; exit 1 if exceeded")
    val.set_defaults(func=_cmd_validate)

    args = parser.parse_args(argv)
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
