"""Tests for gen_lunar_new_year.

Environment variables (tests that need a missing file are skipped):

    GENCAL_EPHEMERIS  path to de440.bsp          (default ~/tommi-data/ephemeris/de440.bsp)
    GENCAL_BUNDLED_LNY reference table to compare with (default: the v1.4.3 table kept in
                       ../../testdata/calendar/lunar-new-year-v1.4.3.json)
    GENCAL_SLOW       set to 1 to run the checks over the whole 1700-2399 range and to check that
                      data/lunar-new-year.json is what the generator writes (about 8 minutes)
"""

from __future__ import annotations

import datetime as dt
import json
import os
from datetime import datetime, timezone
from pathlib import Path

import pytest

import gen_lunar_new_year as g

EPHEMERIS = Path(os.environ.get("GENCAL_EPHEMERIS", Path.home() / "tommi-data/ephemeris/de440.bsp"))
REPO = Path(__file__).resolve().parents[2]
BUNDLED = Path(
    os.environ.get("GENCAL_BUNDLED_LNY") or REPO / "testdata" / "calendar" / "lunar-new-year-v1.4.3.json"
)
CURRENT = REPO / "data" / "lunar-new-year.json"

needs_ephemeris = pytest.mark.skipif(not EPHEMERIS.exists(), reason="de440.bsp not found")
needs_bundled = pytest.mark.skipif(
    not BUNDLED.exists(), reason="reference lunar-new-year.json not found; set GENCAL_BUNDLED_LNY"
)
slow = pytest.mark.skipif(not os.environ.get("GENCAL_SLOW"), reason="set GENCAL_SLOW=1")

CHINA = g.ZONES["china"]


def utc(*parts: int) -> datetime:
    return datetime(*parts, tzinfo=timezone.utc)


def mmdd(new_year: g.LunarNewYear) -> str:
    return new_year.day.strftime("%m-%d")


def test_fixed_zone_offsets_the_date():
    assert g.local_date(utc(2020, 1, 24, 16, 30), g.ZONES["utc+8"]) == dt.date(2020, 1, 25)
    assert g.local_date(utc(2020, 1, 24, 16, 30), g.ZONES["utc+7"]) == dt.date(2020, 1, 24)


def test_china_zone_uses_beijing_mean_time_before_1929():
    before = utc(1916, 2, 3, 16, 10)  # 00:10 at UTC+8, 23:55:40 at Beijing mean time
    assert g.local_date(before, g.ZONES["utc+8"]) == dt.date(1916, 2, 4)
    assert g.local_date(before, CHINA) == dt.date(1916, 2, 3)
    assert CHINA.offset_at(utc(1928, 12, 31, 23, 59)) < 8.0
    assert CHINA.offset_at(utc(1929, 1, 1)) == 8.0


def test_minutes_from_local_midnight_is_symmetric():
    zone = g.ZONES["utc+8"]
    assert g.minutes_from_local_midnight(utc(2020, 1, 24, 16, 5), zone) == pytest.approx(5.0)
    assert g.minutes_from_local_midnight(utc(2020, 1, 24, 15, 55), zone) == pytest.approx(5.0)
    assert g.minutes_from_local_midnight(utc(2020, 1, 24, 4, 0), zone) == pytest.approx(720.0)


def test_inverted_range_is_rejected():
    with pytest.raises(ValueError, match="after last_year"):
        g.compute_lunar_new_years(EPHEMERIS, 2024, 2023, CHINA)


def test_json_document_has_the_bundled_shape():
    new_years = [g.LunarNewYear(2020, dt.date(2020, 1, 25), utc(2020, 1, 24, 21, 42), 6)]
    assert g.to_json_document(new_years) == {"lunarNewYearDates": {"2020": "01-25"}}


@needs_ephemeris
@pytest.mark.parametrize(
    "published",
    [utc(2000, 1, 6, 18, 14), utc(2024, 1, 11, 11, 57), utc(2024, 2, 9, 22, 59)],
)
def test_new_moons_match_published_minutes(published):
    moons = g.find_new_moons(EPHEMERIS, published - dt.timedelta(days=3), published + dt.timedelta(days=3))
    assert len(moons) == 1
    assert abs((moons[0] - published).total_seconds()) <= 60


@needs_ephemeris
def test_new_moons_are_one_lunation_apart():
    moons = g.find_new_moons(EPHEMERIS, utc(2020, 1, 1), utc(2024, 1, 1))
    gaps = [(b - a).total_seconds() / 86400 for a, b in zip(moons, moons[1:])]
    assert len(moons) in (49, 50)
    assert all(29.2 < gap < 29.9 for gap in gaps)


@pytest.fixture(scope="module")
def new_years_2019_2034():
    return {n.year: n for n in g.compute_lunar_new_years(EPHEMERIS, 2019, 2034, CHINA)}


@needs_ephemeris
@pytest.mark.parametrize(
    ("year", "want"),
    [(2019, "02-05"), (2020, "01-25"), (2023, "01-22"), (2024, "02-10"), (2033, "01-31"), (2034, "02-19")],
)
def test_published_lunar_new_year_dates(new_years_2019_2034, year, want):
    assert mmdd(new_years_2019_2034[year]) == want


@needs_ephemeris
def test_leap_month_is_found_after_month_one_in_2020(new_years_2019_2034):
    # 2020 has a leap fourth month: month 11, 12, 1, 2, 3, 4, then the leap month.
    assert new_years_2019_2034[2020].leap_month_index == 6
    assert new_years_2019_2034[2021].leap_month_index is None


@needs_ephemeris
def test_gaps_between_consecutive_years_are_valid(new_years_2019_2034):
    days = [new_years_2019_2034[y].day for y in range(2019, 2035)]
    gaps = [(b - a).days for a, b in zip(days, days[1:])]
    assert all(g_ in (353, 354, 355, 383, 384, 385) for g_ in gaps)


@needs_ephemeris
def test_1916_follows_beijing_mean_time_not_utc8():
    # The new moon is 00:05 at UTC+8 on 4 February but still 3 February at Beijing mean time.
    [china] = g.compute_lunar_new_years(EPHEMERIS, 1916, 1916, CHINA)
    [utc8] = g.compute_lunar_new_years(EPHEMERIS, 1916, 1916, g.ZONES["utc+8"])
    assert (mmdd(china), mmdd(utc8)) == ("02-03", "02-04")


@needs_ephemeris
@needs_bundled
def test_writer_reproduces_the_bundled_file_layout(tmp_path):
    raw = BUNDLED.read_text(encoding="utf-8")
    out = tmp_path / "lunar-new-year.json"
    g.write_json(json.loads(raw), out)
    assert out.read_text(encoding="utf-8") == raw


@needs_ephemeris
@needs_bundled
def test_generator_reproduces_all_bundled_dates():
    bundled = json.loads(BUNDLED.read_text(encoding="utf-8"))["lunarNewYearDates"]
    years = sorted(int(y) for y in bundled)
    new_years = g.compute_lunar_new_years(EPHEMERIS, years[0], years[-1], CHINA)
    assert {str(n.year): mmdd(n) for n in new_years} == bundled


@needs_ephemeris
@slow
def test_first_release_range_has_valid_gaps_and_dates():
    new_years = g.compute_lunar_new_years(EPHEMERIS, 1700, 2399, CHINA)
    assert [n.year for n in new_years] == list(range(1700, 2400))
    gaps = [(b.day - a.day).days for a, b in zip(new_years, new_years[1:])]
    assert set(gaps) <= {353, 354, 355, 383, 384, 385}
    assert all(dt.date(n.year, 1, 21) <= n.day <= dt.date(n.year, 2, 21) for n in new_years)


@needs_ephemeris
@slow
@pytest.mark.skipif(not CURRENT.exists(), reason="data/lunar-new-year.json not found")
def test_data_file_is_what_the_generator_writes(tmp_path):
    current = json.loads(CURRENT.read_text(encoding="utf-8"))["lunarNewYearDates"]
    years = sorted(int(y) for y in current)
    out = tmp_path / "lunar-new-year.json"
    g.write_json(g.to_json_document(g.compute_lunar_new_years(EPHEMERIS, years[0], years[-1], CHINA)), out)
    assert out.read_text(encoding="utf-8") == CURRENT.read_text(encoding="utf-8")
