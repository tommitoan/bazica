package bazica

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/model"
)

// The released range is pinned on purpose: a change of the calendar tables that
// moves it must be a deliberate edit of this test.
func TestSupportedYears(t *testing.T) {
	first, last, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	if first != 1700 || last != 2399 {
		t.Errorf("SupportedYears() = %d, %d, want 1700, 2399", first, last)
	}
}

func TestGetBaziChartAtBothEndsOfTheSupportedRange(t *testing.T) {
	ict := time.FixedZone("ICT", 7*3600)
	first, last, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		date time.Time
	}{
		{"first minute", time.Date(first, 1, 1, 0, 0, 0, 0, ict)},
		{"first noon", time.Date(first, 1, 1, 12, 0, 0, 0, ict)},
		{"last noon", time.Date(last, 12, 31, 12, 0, 0, 0, ict)},
		// From 23:00 the day counts as the next one, so the calendar date is already 1 January of the year after.
		{"last minute", time.Date(last, 12, 31, 23, 59, 0, 0, ict)},
	}
	for _, gender := range []int{model.GenderFemale, model.GenderMale} {
		for _, tc := range tests {
			t.Run(fmt.Sprintf("%s gender %d", tc.name, gender), func(t *testing.T) {
				chart, err := GetBaziChart(tc.date, ict, gender)
				if err != nil {
					t.Fatalf("GetBaziChart(%s) error = %v", tc.date, err)
				}
				if chart.FourPillar.YearPillar == nil || chart.LuckPillars == nil || len(chart.LuckPillars.LuckPillars) == 0 {
					t.Errorf("incomplete chart for %s", tc.date)
				}
			})
		}
	}
}

func TestGetBaziChartJustOutsideTheSupportedRange(t *testing.T) {
	ict := time.FixedZone("ICT", 7*3600)
	first, last, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	for name, date := range map[string]time.Time{
		"last day before the range": time.Date(first-1, 12, 31, 12, 0, 0, 0, ict),
		"first day after the range": time.Date(last+1, 1, 1, 12, 0, 0, 0, ict),
		"middle of the year after":  time.Date(last+1, 6, 1, 12, 0, 0, 0, ict),
		"middle of the year before": time.Date(first-1, 6, 1, 12, 0, 0, 0, ict),
	} {
		if _, err := GetBaziChart(date, ict, model.GenderMale); !errors.Is(err, model.ErrDateOutOfRange) {
			t.Errorf("%s: error = %v, want ErrDateOutOfRange", name, err)
		}
	}
}
