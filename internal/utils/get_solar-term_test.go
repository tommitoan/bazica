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
