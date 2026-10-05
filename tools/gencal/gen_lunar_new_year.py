"""Generate Lunar New Year dates from new moons and the winter solstice.

Rules of the Chinese-derived lunisolar calendar, applied to civil dates in a named zone:

- A month starts on the local date of the new moon (conjunction in ecliptic longitude).
- The month that contains the winter solstice is month 11.
- Between two such months there are 12 or 13 months. With 13, the first month that
  holds no principal term (a multiple of 30 degrees of solar longitude) is the leap
  month and takes the number of the month before it.
- Lunar New Year is the first day of month 1.

Positions come from the same JPL DE440 ephemeris as gen_solar_terms.py, which supplies
the principal terms. The output has the shape of bazica's data/lunar-new-year.json.
"""

from __future__ import annotations

import argparse
import json
import sys
from bisect import bisect_left, bisect_right
from dataclasses import dataclass
from datetime import date, datetime, timedelta, timezone
from pathlib import Path

import numpy as np
from skyfield import almanac, api

import gen_solar_terms as solar

_SEARCH_STEP_DAYS = 7.0  # a quarter of a lunation is about 7.4 days, so no phase is skipped
_EPSILON_DAYS = solar._EPSILON_DAYS
# A phase crossing must have an elongation this close to a quarter (about 4 s of Moon motion).
_MAX_RESIDUAL_DEGREES = 2e-3

# Every multiple of 30 degrees is a principal term: every second name in solar.TERM_NAMES.
_PRINCIPAL_TERMS = frozenset(solar.TERM_NAMES[::2])
_WINTER_SOLSTICE = "winter_solstice"


@dataclass(frozen=True)
class Zone:
    """The civil time a calendar follows: a fixed offset, optionally different before a cutover."""

    name: str
    offset_hours: float
    early_offset_hours: float | None = None
    cutover: datetime | None = None

    def offset_at(self, instant: datetime) -> float:
        """UTC offset in hours that applies at a UTC instant."""
        if self.cutover is not None and self.early_offset_hours is not None and instant < self.cutover:
            return self.early_offset_hours
        return self.offset_hours


# Beijing mean time (116 deg 25 min east, 7 h 45 min 40 s) was the reference of the
# Chinese calendar until standard time UTC+8 was adopted on 1929-01-01.
_BEIJING_MEAN_TIME_HOURS = 7 + 45 / 60 + 40 / 3600
ZONES: dict[str, Zone] = {
    "utc+7": Zone("utc+7", 7.0),
    "utc+8": Zone("utc+8", 8.0),
    "china": Zone("china", 8.0, _BEIJING_MEAN_TIME_HOURS, datetime(1929, 1, 1, tzinfo=timezone.utc)),
}


@dataclass(frozen=True)
class LunarNewYear:
    """Lunar New Year of one solar year, with the facts that decided it."""

    year: int
    day: date
    new_moon_utc: datetime
    leap_month_index: int | None  # months after month 11 of the previous year; None if no leap


def local_time(instant: datetime, zone: Zone) -> datetime:
    """The clock reading of a UTC instant in a zone, as a naive datetime."""
    return (instant + timedelta(hours=zone.offset_at(instant))).replace(tzinfo=None)


def local_date(instant: datetime, zone: Zone) -> date:
    """Civil date of a UTC instant in a zone."""
    return local_time(instant, zone).date()


def minutes_from_local_midnight(instant: datetime, zone: Zone) -> float:
    """Distance of an instant from the nearest local midnight, in minutes."""
    local = local_time(instant, zone)
    since_midnight = local.hour * 60 + local.minute + local.second / 60.0
    return min(since_midnight, 1440.0 - since_midnight)


def find_new_moons(ephemeris_path: Path, start: datetime, end: datetime) -> list[datetime]:
    """Return the UTC instants of every new moon in [start, end), rounded to the millisecond."""
    ts = api.load.timescale()
    eph = api.load_file(str(ephemeris_path))
    earth, sun, moon = eph["earth"], eph["sun"], eph["moon"]

    def elongation(t) -> np.ndarray:
        _, sun_lon, _ = earth.at(t).observe(sun).apparent().ecliptic_latlon(epoch="date")
        _, moon_lon, _ = earth.at(t).observe(moon).apparent().ecliptic_latlon(epoch="date")
        return (moon_lon.degrees - sun_lon.degrees) % 360.0

    def quarter(t) -> np.ndarray:
        return np.floor(elongation(t) / 90.0).astype(int) % 4

    quarter.step_days = _SEARCH_STEP_DAYS

    t0 = ts.from_datetime(start - timedelta(days=10))
    t1 = ts.from_datetime(end + timedelta(days=10))
    times, _ = almanac.find_discrete(t0, t1, quarter, epsilon=_EPSILON_DAYS)

    # See gen_solar_terms.find_solar_terms: the step function can flicker at float
    # resolution, so merge repeats and decide what each crossing is from the elongation.
    seconds = times.tt * 86400.0
    keep = np.concatenate(([True], np.diff(seconds) > solar._MIN_SEPARATION_SECONDS))
    times = times[keep]

    steps = elongation(times) / 90.0
    nearest = np.rint(steps)
    residual_degrees = np.abs(steps - nearest) * 90.0
    if residual_degrees.size and residual_degrees.max() > _MAX_RESIDUAL_DEGREES:
        raise ValueError(f"a phase crossing is {residual_degrees.max():.6f} degrees from a quarter")

    moons: list[datetime] = []
    for t, quarter_index in zip(times, nearest.astype(int) % 4):
        instant = solar._round_to_millisecond(t.utc_datetime())
        if quarter_index == 0 and start <= instant < end:
            moons.append(instant)
    return moons


def compute_lunar_new_years(
    ephemeris_path: Path, first_year: int, last_year: int, zone: Zone
) -> list[LunarNewYear]:
    """Lunar New Year for every solar year in first_year..last_year in a zone."""
    if first_year > last_year:
        raise ValueError(f"first_year {first_year} is after last_year {last_year}")

    terms = solar.find_solar_terms(ephemeris_path, first_year - 1, last_year)
    solstices = {t.instant.year: t.instant for t in terms if t.name == _WINTER_SOLSTICE}
    principal_days = sorted(local_date(t.instant, zone) for t in terms if t.name in _PRINCIPAL_TERMS)

    window_start = datetime(first_year - 1, 11, 1, tzinfo=timezone.utc)
    window_end = datetime(last_year + 1, 1, 15, tzinfo=timezone.utc)
    new_moons = find_new_moons(ephemeris_path, window_start, window_end)
    moon_days = [local_date(m, zone) for m in new_moons]

    def has_principal_term(month: int) -> bool:
        first = bisect_left(principal_days, moon_days[month])
        return first < len(principal_days) and principal_days[first] < moon_days[month + 1]

    results: list[LunarNewYear] = []
    for year in range(first_year, last_year + 1):
        previous = bisect_right(moon_days, local_date(solstices[year - 1], zone)) - 1
        current = bisect_right(moon_days, local_date(solstices[year], zone)) - 1
        months = current - previous

        leap_index: int | None = None
        if months == 13:
            for month in range(previous + 1, current):
                if not has_principal_term(month):
                    leap_index = month - previous
                    break
            if leap_index is None:
                raise ValueError(f"{year}: 13 months but every month has a principal term")
        elif months != 12:
            raise ValueError(f"{year}: {months} months between winter solstices")

        # Month 11 is `previous`; month 1 follows month 12, and a leap 11 or 12 pushes it one later.
        first_month = previous + (3 if leap_index in (1, 2) else 2)
        day = moon_days[first_month]
        if day.year != year:
            raise ValueError(f"{year}: Lunar New Year computed as {day}")
        results.append(LunarNewYear(year, day, new_moons[first_month], leap_index))
    return results


def to_json_document(new_years: list[LunarNewYear]) -> dict:
    """Build the lunar-new-year.json structure."""
    return {"lunarNewYearDates": {str(n.year): n.day.strftime("%m-%d") for n in new_years}}


def write_json(document: dict, out_path: Path) -> None:
    """Write the document in the layout of the bundled file (2-space indent, trailing newline)."""
    out_path.write_text(json.dumps(document, indent=2, ensure_ascii=False) + "\n", encoding="utf-8")


def _near_midnight(new_years: list[LunarNewYear], zone: Zone, minutes: float) -> list[tuple[int, float]]:
    pairs = [(n.year, minutes_from_local_midnight(n.new_moon_utc, zone)) for n in new_years]
    return [(year, round(distance, 2)) for year, distance in pairs if distance < minutes]


def _cmd_generate(args: argparse.Namespace) -> int:
    zone = ZONES[args.zone]
    new_years = compute_lunar_new_years(args.ephemeris, args.first_year, args.last_year, zone)
    write_json(to_json_document(new_years), args.out)
    print(f"wrote {len(new_years)} years ({zone.name}) to {args.out}")
    near = _near_midnight(new_years, zone, args.near_minutes)
    print(f"new moon of Lunar New Year within {args.near_minutes:g} min of local midnight: {near}")
    return 0


def _cmd_validate(args: argparse.Namespace) -> int:
    bundled = json.loads(args.bundled.read_text(encoding="utf-8"))["lunarNewYearDates"]
    years = sorted(int(y) for y in bundled)
    reproduced_by_first_zone = True
    for index, name in enumerate(args.zones):
        zone = ZONES[name]
        new_years = compute_lunar_new_years(args.ephemeris, years[0], years[-1], zone)
        mismatches = [n for n in new_years if n.day.strftime("%m-%d") != bundled[str(n.year)]]
        print(f"{zone.name}: {len(new_years) - len(mismatches)} of {len(new_years)} dates reproduced")
        for n in mismatches:
            print(
                f"  {n.year}: generated {n.day:%m-%d}, bundled {bundled[str(n.year)]}; "
                f"new moon {local_time(n.new_moon_utc, zone):%Y-%m-%d %H:%M:%S} local, "
                f"leap index {n.leap_month_index}"
            )
        print(f"  within {args.near_minutes:g} min of local midnight: "
              f"{_near_midnight(new_years, zone, args.near_minutes)}")
        if index == 0 and mismatches:
            reproduced_by_first_zone = False
    return 0 if reproduced_by_first_zone else 1


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="command", required=True)

    gen = sub.add_parser("generate", help="write lunar-new-year.json for a range of years")
    gen.add_argument("--ephemeris", type=Path, required=True, help="path to de440.bsp")
    gen.add_argument("--first-year", type=int, required=True)
    gen.add_argument("--last-year", type=int, required=True)
    gen.add_argument("--zone", choices=sorted(ZONES), default="china")
    gen.add_argument("--near-minutes", type=float, default=15.0)
    gen.add_argument("--out", type=Path, required=True)
    gen.set_defaults(func=_cmd_generate)

    val = sub.add_parser("validate", help="compare with a bundled lunar-new-year.json; exit 1 if the first zone differs")
    val.add_argument("--ephemeris", type=Path, required=True, help="path to de440.bsp")
    val.add_argument("--bundled", type=Path, required=True, help="path to data/lunar-new-year.json")
    val.add_argument("--zones", nargs="+", choices=sorted(ZONES), default=["china", "utc+8", "utc+7"])
    val.add_argument("--near-minutes", type=float, default=15.0)
    val.set_defaults(func=_cmd_validate)

    args = parser.parse_args(argv)
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
