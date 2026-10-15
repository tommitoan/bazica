package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/model"
)

func TestGetSolarTerm(t *testing.T) {
	utc := time.UTC
	tests := []struct {
		name          string
		date          time.Time
		wantTerm      string
		wantPassedMin int // minimum minutes since the preceding initial term
		wantRemainMin int // minimum minutes until the next initial term
	}{
		{"mid month initial term", time.Date(2024, 6, 10, 0, 0, 0, 0, utc), model.GrainInEar, 0, 0},
		{"midpoint term", time.Date(2024, 6, 25, 0, 0, 0, 0, utc), model.SummerSolstice, 1, 1},
		{"after winter solstice", time.Date(2025, 12, 25, 12, 0, 0, 0, utc), model.WinterSolstice, 1, 1},
		{"last day of year", time.Date(2024, 12, 31, 23, 0, 0, 0, utc), model.WinterSolstice, 1, 1},
		{"between major snow and winter solstice", time.Date(2024, 12, 10, 0, 0, 0, 0, utc), model.MajorSnow, 1, 1},
		{"before first term of year", time.Date(2024, 1, 3, 12, 0, 0, 0, utc), model.MajorSnow, 1, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			term, passed, remaining, err := GetSolarTerm(tc.date)
			if err != nil {
				t.Fatalf("GetSolarTerm() error = %v", err)
			}
			if term != tc.wantTerm {
				t.Errorf("term = %s, want %s", term, tc.wantTerm)
			}
			if passed < tc.wantPassedMin || remaining < tc.wantRemainMin {
				t.Errorf("passed/remaining = %d/%d, want at least %d/%d", passed, remaining, tc.wantPassedMin, tc.wantRemainMin)
			}
		})
	}
}

func TestGetSolarTermWindow(t *testing.T) {
	loc := time.FixedZone("Vietnam", 7*60*60)
	cases := []struct {
		name, want, next string
		birth            time.Time
	}{
		{"Grain in Ear", model.GrainInEar, model.SummerSolstice, time.Date(1995, 6, 8, 11, 15, 0, 0, loc)},
		{"before first term", model.WinterSolstice, model.MinorCold, time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC)},
		{"end of year", model.WinterSolstice, model.MinorCold, time.Date(2024, 12, 31, 12, 0, 0, 0, time.UTC)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			window, err := GetSolarTermWindow(tc.birth)
			if err != nil {
				t.Fatal(err)
			}
			if window.Name != tc.want || window.NextName != tc.next {
				t.Errorf("terms = %s/%s, want %s/%s", window.Name, window.NextName, tc.want, tc.next)
			}
			if tc.birth.Before(window.StartedAt) || !tc.birth.Before(window.NextAt) {
				t.Errorf("birth %s outside [%s, %s)", tc.birth, window.StartedAt, window.NextAt)
			}
			if window.StartedAt.Location() != tc.birth.Location() || window.NextAt.Location() != tc.birth.Location() {
				t.Error("term times do not use the birth location")
			}
		})
	}
}

func TestGetSolarTermWindowAtEveryBoundary(t *testing.T) {
	years, err := loadSolarTerms()
	if err != nil {
		t.Fatal(err)
	}
	terms := years[2024]
	if len(terms) != 24 {
		t.Fatalf("2024 has %d terms, want 24", len(terms))
	}
	for _, term := range terms {
		before, err := GetSolarTermWindow(term.time.Add(-time.Millisecond))
		if err != nil {
			t.Fatal(err)
		}
		at, err := GetSolarTermWindow(term.time)
		if err != nil {
			t.Fatal(err)
		}
		after, err := GetSolarTermWindow(term.time.Add(time.Millisecond))
		if err != nil {
			t.Fatal(err)
		}
		if before.NextName != term.name || at.Name != term.name || after.Name != term.name {
			t.Errorf("%s boundary: before=%s/%s at=%s after=%s", term.name, before.Name, before.NextName, at.Name, after.Name)
		}
	}
}

// The two durations around a date must add up to the length of the Ba-zi month.
func TestGetSolarTermMonthLength(t *testing.T) {
	for _, d := range []time.Time{
		time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
	} {
		_, passed, remaining, err := GetSolarTerm(d)
		if err != nil {
			t.Fatalf("%s: %v", d.Format("2006-01-02"), err)
		}
		days := float64(passed+remaining) / (60 * 24)
		if days < 28 || days > 33 {
			t.Errorf("%s: month length %.1f days outside 28-33", d.Format("2006-01-02"), days)
		}
	}
}

func TestGetSolarTermErrors(t *testing.T) {
	tests := []struct {
		name    string
		date    time.Time
		wantErr error
	}{
		{"year without data", time.Date(2401, 6, 1, 0, 0, 0, 0, time.UTC), model.ErrDateOutOfRange},
		{"next year without data", time.Date(2400, 6, 1, 0, 0, 0, 0, time.UTC), model.ErrDateOutOfRange},
		{"previous year without data", time.Date(1699, 6, 1, 0, 0, 0, 0, time.UTC), model.ErrDateOutOfRange},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, _, err := GetSolarTerm(tc.date); !errors.Is(err, tc.wantErr) {
				t.Errorf("error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// Lichun is Start of Spring in early February; the table covers the years the library supports.
func TestStartOfSpring(t *testing.T) {
	got, err := StartOfSpring(2026)
	if err != nil {
		t.Fatal(err)
	}
	if got.UTC().Format("2006-01-02") != "2026-02-03" && got.UTC().Format("2006-01-02") != "2026-02-04" {
		t.Errorf("Lichun 2026 = %s, want 3 or 4 February", got.UTC())
	}
	if first, last, err := SolarTermYears(); err != nil {
		t.Fatal(err)
	} else {
		for _, year := range []int{first - 1, last + 1, 1500, 3000} {
			if _, err := StartOfSpring(year); !errors.Is(err, model.ErrDateOutOfRange) {
				t.Errorf("StartOfSpring(%d) error = %v, want ErrDateOutOfRange", year, err)
			}
		}
	}
}
