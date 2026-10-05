# gencal — calendar table generator for bazica

Generates bazica's two calendar tables from the JPL DE440 ephemeris through `skyfield`:

- `gen_solar_terms.py` writes `solar-term.json` (the 24 solar terms per year).
- `gen_lunar_new_year.py` writes `lunar-new-year.json` (the date of Lunar New Year per year).

Both write the layout of the files in `data/`, byte for byte (see the layout tests). `data/solar-term.json`
(1699-2400) and `data/lunar-new-year.json` (1700-2399, zone `china`) are exactly what these tools write; the
tables of v1.4.3 that they replaced are kept in `testdata/calendar/` for comparison.

A solar term is the instant the Sun's **apparent ecliptic longitude of date** reaches a
multiple of 15 degrees. Times are UTC, rounded to the millisecond. Delta T comes from
skyfield's built-in timescale (version 1.55); the value used at the start of each year is
listed under "Accuracy" below.

## Setup

Python 3.10 or newer. The ephemeris file is not part of any repository.

Identical on macOS, Ubuntu and Fedora:

```bash
python3 -m venv ~/tommi-data/gencal-venv
~/tommi-data/gencal-venv/bin/pip install -r requirements.txt
mkdir -p ~/tommi-data/ephemeris
curl -L -o ~/tommi-data/ephemeris/de440.bsp \
  https://naif.jpl.nasa.gov/pub/naif/generic_kernels/spk/planets/de440.bsp
```

`de440.bsp` is 119,799,808 bytes (about 114 MiB) and covers 1550 to 2650.
Its SHA-256 on 2026-10-05: `a4ce9bf9b3282becc9f4b2ac3cebe03a2ae7599981aabd7265fd8482fff7c4b5`.

OS-specific, only if the venv module is missing:

- Ubuntu only: `sudo apt install python3-venv`
- Fedora only: `python3 -m venv` ships with `python3`; nothing to install
- macOS only: use Python from python.org or Homebrew (`brew install python`)

## Usage (solar terms)

Compare the generator with the table of v1.4.3 (the current `data/` file would compare with itself):

```bash
~/tommi-data/gencal-venv/bin/python gen_solar_terms.py validate \
  --ephemeris ~/tommi-data/ephemeris/de440.bsp \
  --bundled ../../testdata/calendar/solar-term-v1.4.3.json
```

Generate a range (one margin year on each side of the supported range is the caller's choice):

```bash
~/tommi-data/gencal-venv/bin/python gen_solar_terms.py generate \
  --ephemeris ~/tommi-data/ephemeris/de440.bsp \
  --first-year 1699 --last-year 2400 --out solar-term.json
```

1699-2400 takes about one minute and gives 702 years, 969 KB.

Tests (the comparison tests read the v1.4.3 tables in `../../testdata/calendar/`; set `GENCAL_BUNDLED` to use another file,
`GENCAL_BUNDLED_LNY` for `lunar-new-year.json`, `GENCAL_EPHEMERIS` to use another ephemeris path; tests that need a missing file are skipped):

```bash
~/tommi-data/gencal-venv/bin/python -m pytest -q
```

This folder is a Python tool inside a Go repository: it is not part of the Go module and adds no dependency.

## Result against the v1.4.3 solar-term table (4,848 terms, 1899-2100)

| Period | Mean signed difference (generated - bundled) | Largest absolute |
|---|---|---|
| 1900-1919 | -36.3 s | 49.8 s |
| 1920-1939 | -21.9 s | 24.9 s |
| 1940-1959 | -15.6 s | 20.7 s |
| 1960-1979 | -4.3 s | 10.5 s |
| 1980-2039 | -1.2 s to -0.5 s (std 0.03 s in 1980-2019) | 1.2 s |
| 2040-2059 | +0.9 s | 8.5 s |
| 2060-2079 | +2.4 s | 9.3 s |
| 2080-2100 | +4.0 s | 4.9 s |

Median difference over all terms: 2.7 s. The difference is smooth in time, so it is a Delta T
model difference, not noise: the bundled table follows a Delta T that agrees with skyfield's within
about a second from 1980 and drifts away earlier and later.

Rows of the v1.4.3 file that disagreed by more (the table above leaves them out). The regenerated
`data/solar-term.json` replaces all of them:

- `winter_solstice` 1903 is `1903-12-23 00:00:03.866`; the generator gives `00:19:45.238`
  (19 minutes 41 seconds off).
- `minor_cold` of 17 years from 2041 on is off by 17 to 98 seconds (2053: -98 s, 2077: -90 s,
  2081: +73 s). Earlier these rows held a December instant; the repair is close but approximate.
- The margin year 1899 is up to 3.6 hours off for its first 22 terms (it is only used for births in
  January 1900, which need the December terms of 1899).

## Lunar New Year (`gen_lunar_new_year.py`)

Rules, applied to civil dates in a named zone:

- A month starts on the local date of the new moon (conjunction in ecliptic longitude).
- The month that contains the winter solstice is month 11.
- Between two such months there are 12 or 13 months. With 13, the first month with no principal
  term (a multiple of 30 degrees of solar longitude) is the leap month.
- Lunar New Year is the first day of month 1. A leap month 11 or 12 pushes it one month later.

```bash
~/tommi-data/gencal-venv/bin/python gen_lunar_new_year.py validate \
  --ephemeris ~/tommi-data/ephemeris/de440.bsp --bundled ../../testdata/calendar/lunar-new-year-v1.4.3.json

~/tommi-data/gencal-venv/bin/python gen_lunar_new_year.py generate \
  --ephemeris ~/tommi-data/ephemeris/de440.bsp \
  --first-year 1700 --last-year 2399 --zone china --out lunar-new-year.json
```

Validation takes about one minute and generating 1700-2399 about four. The full-range test only runs
with `GENCAL_SLOW=1`.

### Which time zone the bundled table follows

| Zone | Dates reproduced (1900-2099) |
|---|---|
| `china`: UTC+8 from 1929-01-01, Beijing mean time (UTC+7:45:40) before | **200 of 200** |
| `utc+8` | 199 of 200 (1916: new moon 00:05 on 4 February at UTC+8, 23:50 on 3 February at Beijing mean time) |
| `utc+7` (Vietnam) | 191 of 200 (1903, 1935, 1965, 1968, 1969, 1985, 2007, 2030, 2053) |

The v1.4.3 table therefore follows the Chinese calendar convention, not the Vietnamese one. The regenerated
table keeps every date of 1900-2099, so the generator uses `china`. A Vietnamese variant (UTC+7) would change nine of
those years and is out of scope.

Only one year (1916) separates `china` from `utc+8` inside 1900-2099, so the pre-1929 Beijing mean time
rule rests on that single case there. For 1700-2399 the two zones also differ in 1896 only
(`china` 13 February, `utc+8` 14 February). Phase 4 cross-checks these against an independent source.

### Range 1700-2399 (zone `china`)

All 700 dates have a gap of 353-355 or 383-385 days to the next one, and fall between 21 January and
21 February (the Go integrity test for 1900-2099 stops at 20 February; widen it with the data). Sample:
1700 `02-19`, 1800 `01-25`, 2100 `02-09`, 2399 `02-07`.

### Years that are uncertain

When the new moon of Lunar New Year falls close to local midnight, a small error in the time moves the
date. Distance from midnight, in minutes, for the years within 15 minutes:

| Period | Years (minutes from local midnight) |
|---|---|
| 1700-1899 | 1808 (7.9), 1893 (1.5), 1896 (2.6) |
| 1900-2099 | 1916 (9.5), 1954 (4.8), 1966 (14.0), 1988 (5.8), 2007 (14.3), 2027 (3.9), 2030 (7.5) |
| 2100-2399 | 2215 (10.1), 2261 (8.1), 2299 (2.1), 2303 (8.5), 2333 (0.27), 2375 (11.2) |

Years up to 2099 that appear here agree with the v1.4.3 table. Past 2100 Delta T is a prediction (see
above), so 2299 and 2333 in particular can land on the other side of midnight; treat the date of those
years as uncertain.

## Accuracy beyond the bundled range

Delta T is a measurement up to the present and a prediction after it. skyfield gives
about 96 s in 2100, 222 s in 2200 and 753 s in 2399; the prediction can be wrong by
minutes at the far end, and a term moves one second for each second of Delta T. Before
1700 and before 1900 the value is a reconstruction with an uncertainty of a second or two.
A birth within that margin of a term is uncertain.

## Implementation notes

- A 10-day search step cannot skip a crossing: two 15-degree crossings are at least 14.7 days apart.
- `find_discrete` epsilon must stay above the float64 resolution of a Julian date (about 40
  microseconds). Asking for 0.1 microsecond did not terminate and used 22 GB for one year.
- At float resolution the step function can flicker across a crossing, so one crossing came
  back three times (2309 autumn equinox) and once with the previous term's value.
  The generator merges crossings closer than a second and names each from the longitude at
  the instant. `test_first_release_range_has_24_distinct_terms_in_every_year` guards this.
