package bazica

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/tommitoan/bazica/model"
)

func referenceModelChart(t *testing.T) *model.BaziChart {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	chart, err := GetBaziChart(time.Date(1995, 6, 8, 11, 1, 0, 0, loc), loc, 1)
	if err != nil {
		t.Fatalf("GetBaziChart() error = %v", err)
	}
	return chart
}

func TestGetAnnualPillarsMatchesTheSharedFixture(t *testing.T) {
	annual, err := GetAnnualPillars(referenceModelChart(t), 1995, 60)
	if err != nil {
		t.Fatalf("GetAnnualPillars() error = %v", err)
	}
	if !reflect.DeepEqual(roundTrip(t, annual), loadFixture(t, "annual_1995_expected.json")) {
		t.Errorf("annual pillars differ from the fixture:\n%s", firstDifference(mustJSON(t, roundTrip(t, annual)), mustJSON(t, loadFixture(t, "annual_1995_expected.json"))))
	}
}

func TestGetAnnualPillarsRows(t *testing.T) {
	annual, err := GetAnnualPillars(referenceModelChart(t), 2026, 3)
	if err != nil {
		t.Fatal(err)
	}
	if annual.FromYear != 2026 || annual.Count != 3 || len(annual.AnnualPillars) != 3 {
		t.Fatalf("got %+v", annual)
	}
	first := annual.AnnualPillars[0]
	// 2026 is Bing-Wu; the chart's Day Master is Geng and its luck pillar is number 4.
	if first.HeavenlyStem.Spelling != "bing" || first.EarthlyBranch.Spelling != "wu" || first.NominalAge != 32 {
		t.Errorf("2026 = %s/%s age %d", first.HeavenlyStem.Spelling, first.EarthlyBranch.Spelling, first.NominalAge)
	}
	if first.TenGod.Code != "ten_god.seven_killings" || first.LuckPillarNumber == nil || *first.LuckPillarNumber != 4 {
		t.Errorf("2026 ten god %s, luck pillar %v", first.TenGod.Code, first.LuckPillarNumber)
	}
}

func TestGetAnnualPillarsLuckPillarNumberIsNullPastTheLastLuckPillar(t *testing.T) {
	annual, err := GetAnnualPillars(referenceModelChart(t), 2099, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got := annual.AnnualPillars[0].LuckPillarNumber; got == nil || *got != 11 {
		t.Errorf("2099 luck pillar = %v, want 11", got)
	}
	data, _ := json.Marshal(model.AnnualPillar{})
	var m map[string]any
	_ = json.Unmarshal(data, &m)
	if v, ok := m["luck_pillar_number"]; !ok || v != nil {
		t.Errorf("empty luck pillar number should serialise as null, got %v (present %v)", v, ok)
	}
}

func TestGetAnnualPillarsRejectsInvalidRanges(t *testing.T) {
	chart := referenceModelChart(t)
	for name, tc := range map[string]struct {
		chart     *model.BaziChart
		from, num int
	}{
		"nil chart":            {nil, 1995, 5},
		"zero count":           {chart, 1995, 0},
		"negative count":       {chart, 1995, -3},
		"before year pillar":   {chart, 1994, 5},
		"before supported":     {chart, 1899, 5},
		"past supported range": {chart, 2095, 10},
	} {
		if _, err := GetAnnualPillars(tc.chart, tc.from, tc.num); !errors.Is(err, model.ErrInvalidYearRange) {
			t.Errorf("%s: error = %v, want ErrInvalidYearRange", name, err)
		}
	}
	if _, err := GetAnnualPillars(chart, 2090, 10); err != nil {
		t.Errorf("2090 for 10 years ends at 2099 and must be accepted: %v", err)
	}
}

// A January birth before the Lunar New Year belongs to the previous year pillar,
// so nominal age counts from that year and 1999 is the first valid year.
func TestGetAnnualPillarsUsesTheYearPillarYear(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	chart, err := GetBaziChart(time.Date(2000, 1, 15, 14, 30, 0, 0, loc), loc, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := GetAnnualPillars(chart, 1998, 2); !errors.Is(err, model.ErrInvalidYearRange) {
		t.Errorf("1998 precedes the year pillar: error = %v", err)
	}
	annual, err := GetAnnualPillars(chart, 1999, 2)
	if err != nil {
		t.Fatal(err)
	}
	if annual.AnnualPillars[0].NominalAge != 1 || annual.AnnualPillars[1].NominalAge != 2 {
		t.Errorf("nominal ages = %d, %d, want 1, 2", annual.AnnualPillars[0].NominalAge, annual.AnnualPillars[1].NominalAge)
	}
}
