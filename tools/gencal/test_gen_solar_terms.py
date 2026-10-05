"""Tests for gen_solar_terms.

The ephemeris and the bundled table are looked up through environment variables
so the suite runs on any machine; tests that need them are skipped when missing.

    GENCAL_EPHEMERIS  path to de440.bsp          (default ~/tommi-data/ephemeris/de440.bsp)
    GENCAL_BUNDLED    path to data/solar-term.json of bazica (default: ../../data/solar-term.json
                      when this folder is tools/gencal inside the bazica repository)
"""

from __future__ import annotations

import json
import os
from datetime import datetime, timezone
from pathlib import Path

import pytest

import gen_solar_terms as g

EPHEMERIS = Path(os.environ.get("GENCAL_EPHEMERIS", Path.home() / "tommi-data/ephemeris/de440.bsp"))
BUNDLED = Path(
    os.environ.get("GENCAL_BUNDLED") or Path(__file__).resolve().parents[2] / "data" / "solar-term.json"
)

needs_ephemeris = pytest.mark.skipif(not EPHEMERIS.exists(), reason="de440.bsp not found")
needs_bundled = pytest.mark.skipif(
    not BUNDLED.exists(), reason="bundled solar-term.json not found; set GENCAL_BUNDLED"
)


def utc(*parts: int) -> datetime:
    return datetime(*parts, tzinfo=timezone.utc)


@pytest.fixture(scope="module")
def terms_2020_2030():
    return g.find_solar_terms(EPHEMERIS, 2020, 2030)


def test_term_names_are_24_unique():
    assert len(g.TERM_NAMES) == 24
    assert len(set(g.TERM_NAMES)) == 24


def test_timestamp_round_trip_keeps_milliseconds():
    instant = utc(2023, 1, 5, 15, 4, 50).replace(microsecond=647000)
    text = g.format_timestamp(instant)
    assert text == "2023-01-05 15:04:50.647+00:00"
    assert g.parse_timestamp(text) == instant


def test_format_timestamp_pads_milliseconds():
    assert g.format_timestamp(utc(2000, 1, 1, 0, 0, 0).replace(microsecond=5000)) == (
        "2000-01-01 00:00:00.005+00:00"
    )


def test_inverted_range_is_rejected():
    with pytest.raises(ValueError, match="after last_year"):
        g.find_solar_terms(EPHEMERIS, 2024, 2023)


def test_incomplete_year_is_rejected():
    partial = [g.SolarTerm("minor_cold", utc(2023, 1, 5))]
    with pytest.raises(ValueError, match="has 1 terms"):
        g.terms_by_year(partial)


def test_duplicate_term_in_a_year_is_rejected():
    twice = [g.SolarTerm("minor_cold", utc(2023, 1, 5)), g.SolarTerm("minor_cold", utc(2023, 12, 5))]
    with pytest.raises(ValueError, match="occurs twice"):
        g.terms_by_year(twice)


@needs_ephemeris
def test_every_year_has_24_terms(terms_2020_2030):
    grouped = g.terms_by_year(terms_2020_2030)
    assert sorted(grouped) == list(range(2020, 2031))
    assert all(len(year_terms) == 24 for year_terms in grouped.values())


@needs_ephemeris
def test_terms_are_chronological_and_follow_the_longitude_cycle(terms_2020_2030):
    instants = [t.instant for t in terms_2020_2030]
    assert instants == sorted(instants)
    indexes = [g.TERM_NAMES.index(t.name) for t in terms_2020_2030]
    for previous, current in zip(indexes, indexes[1:]):
        assert current == (previous + 1) % 24
    gaps_days = [(b - a).total_seconds() / 86400 for a, b in zip(instants, instants[1:])]
    assert all(14.5 < gap < 16.5 for gap in gaps_days)


@needs_ephemeris
@pytest.mark.parametrize(
    ("year", "name", "published"),
    [
        (2000, "spring_equinox", utc(2000, 3, 20, 7, 35)),
        (2023, "spring_equinox", utc(2023, 3, 20, 21, 24)),
        (2023, "summer_solstice", utc(2023, 6, 21, 14, 58)),
        (2023, "autumn_equinox", utc(2023, 9, 23, 6, 50)),
        (2023, "winter_solstice", utc(2023, 12, 22, 3, 27)),
        (2024, "spring_equinox", utc(2024, 3, 20, 3, 6)),
    ],
)
def test_equinoxes_and_solstices_match_published_minutes(year, name, published):
    grouped = g.terms_by_year(g.find_solar_terms(EPHEMERIS, year, year))
    assert abs((grouped[year][name] - published).total_seconds()) <= 60


@needs_ephemeris
@needs_bundled
def test_writer_reproduces_the_bundled_file_layout(tmp_path):
    raw = BUNDLED.read_text(encoding="utf-8")
    out = tmp_path / "solar-term.json"
    g.write_json(json.loads(raw), out)
    assert out.read_text(encoding="utf-8") == raw


@needs_ephemeris
@needs_bundled
def test_generator_reproduces_the_bundled_table():
    bundled = json.loads(BUNDLED.read_text(encoding="utf-8"))
    grouped = g.terms_by_year(g.find_solar_terms(EPHEMERIS, 1899, 2100))
    differences, years = g.compare_with_bundled(grouped, bundled)
    assert years[0] == 1899 and years[-1] == 2100

    # 1980-2040: the bundled table tracks the same delta T within about a second.
    modern = [d for d in differences if 1980 <= d.year <= 2040]
    assert max(abs(d.seconds) for d in modern) <= 2.0

    # 1900-2100: every term agrees within a minute, except rows of the bundled
    # table that are known to differ more: the 1903 winter solstice (about 20
    # minutes early) and 17 minor_cold rows from 2041 on, whose repaired values are
    # off by up to about 100 seconds. The margin year 1899 is not checked.
    sizeable = [d for d in differences if d.year >= 1900 and abs(d.seconds) > 60]
    others = [(d.year, d.name) for d in sizeable if d.name != "minor_cold"]
    late_minor_cold = [d for d in sizeable if d.name == "minor_cold"]
    assert others == [(1903, "winter_solstice")]
    assert all(d.year >= 2041 and abs(d.seconds) <= 120 for d in late_minor_cold)


@needs_ephemeris
def test_first_release_range_has_24_distinct_terms_in_every_year():
    # Regression: over a long search the crossing at an autumn equinox was once
    # reported with the previous term's name, so a year held white_dew twice.
    grouped = g.terms_by_year(g.find_solar_terms(EPHEMERIS, 1699, 2400))
    assert sorted(grouped) == list(range(1699, 2401))
    assert all(set(year_terms) == set(g.TERM_NAMES) for year_terms in grouped.values())
