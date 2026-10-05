# gencal — calendar table generator for bazica

Generates `solar-term.json` (the 24 solar terms per year) from the JPL DE440
ephemeris through `skyfield`. The output has the same layout as bazica's
`data/solar-term.json`, byte for byte (see `test_writer_reproduces_the_bundled_file_layout`).

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

## Usage

Reproduce and compare with a bundled table:

```bash
~/tommi-data/gencal-venv/bin/python gen_solar_terms.py validate \
  --ephemeris ~/tommi-data/ephemeris/de440.bsp \
  --bundled ../../data/solar-term.json
```

Generate a range (one margin year on each side of the supported range is the caller's choice):

```bash
~/tommi-data/gencal-venv/bin/python gen_solar_terms.py generate \
  --ephemeris ~/tommi-data/ephemeris/de440.bsp \
  --first-year 1699 --last-year 2400 --out solar-term.json
```

1699-2400 takes about one minute and gives 702 years, 969 KB.

Tests (the table tests read `../../data/solar-term.json`; set `GENCAL_BUNDLED` to use another file,
`GENCAL_EPHEMERIS` to use another ephemeris path; tests that need a missing file are skipped):

```bash
~/tommi-data/gencal-venv/bin/python -m pytest -q
```

This folder is a Python tool inside a Go repository: it is not part of the Go module and adds no dependency.

## Result against `data/solar-term.json` of v1.4.3 (4,848 terms, 1899-2100)

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

Rows of the bundled file that disagree by more (the table above leaves them out):

- `winter_solstice` 1903 is `1903-12-23 00:00:03.866`; the generator gives `00:19:45.238`
  (19 minutes 41 seconds off).
- `minor_cold` of 17 years from 2041 on is off by 17 to 98 seconds (2053: -98 s, 2077: -90 s,
  2081: +73 s). Earlier these rows held a December instant; the repair is close but approximate.
- The margin year 1899 is up to 3.6 hours off for its first 22 terms (it is only used for births in
  January 1900, which need the December terms of 1899).

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
