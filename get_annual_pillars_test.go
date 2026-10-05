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
		"past supported range": {chart, 9995, 10},
		"huge count":           {chart, 1995, 100000},
	} {
		if _, err := GetAnnualPillars(tc.chart, tc.from, tc.num); !errors.Is(err, model.ErrInvalidYearRange) {
			t.Errorf("%s: error = %v, want ErrInvalidYearRange", name, err)
		}
	}
	if _, err := GetAnnualPillars(chart, 2090, 10); err != nil {
		t.Errorf("2090 for 10 years ends at 2099 and must be accepted: %v", err)
	}
	if _, err := GetAnnualPillars(chart, 2095, 10); err != nil {
		t.Errorf("2095 for 10 years ends at 2104; the table is not limited to 2099: %v", err)
	}
	if _, err := GetAnnualPillars(chart, 9999, 1); err != nil {
		t.Errorf("9999 is the last supported year: %v", err)
	}
}

// The table continues past 2099 with the sixty-year cycle, the nominal age and the luck pillars of the chart.
func TestGetAnnualPillarsContinuesPast2099(t *testing.T) {
	chart := referenceModelChart(t) // luck pillar 11 of this chart runs 2096 to 2105
	annual, err := GetAnnualPillars(chart, 2095, 25)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(annual.AnnualPillars); got != 25 || annual.AnnualPillars[24].Year != 2119 {
		t.Fatalf("rows = %d, last year %d", got, annual.AnnualPillars[24].Year)
	}
	byYear := map[int]model.AnnualPillar{}
	for _, row := range annual.AnnualPillars {
		byYear[row.Year] = row
	}
	// 2100 is a Geng-Shen year (the same pillar as 2040); the nominal age keeps counting from 1995.
	if r := byYear[2100]; r.HeavenlyStem.Spelling != "geng" || r.EarthlyBranch.Spelling != "shen" || r.NominalAge != 106 {
		t.Errorf("2100 = %s/%s age %d", r.HeavenlyStem.Spelling, r.EarthlyBranch.Spelling, r.NominalAge)
	}
	// The same pillar as sixty years earlier, in every row.
	earlier, err := GetAnnualPillars(chart, 2035, 25)
	if err != nil {
		t.Fatal(err)
	}
	for i, row := range annual.AnnualPillars {
		old := earlier.AnnualPillars[i]
		if row.GanZhi.Name != old.GanZhi.Name || row.HeavenlyStem.Value != old.HeavenlyStem.Value || row.EarthlyBranch.Value != old.EarthlyBranch.Value || row.TenGod != old.TenGod {
			t.Errorf("%d differs from %d: %+v vs %+v", row.Year, old.Year, row.GanZhi, old.GanZhi)
		}
	}
	// Luck pillar 11 ends in 2105; after it the number is null.
	if got := byYear[2105].LuckPillarNumber; got == nil || *got != 11 {
		t.Errorf("2105 luck pillar = %v, want 11", got)
	}
	if got := byYear[2106].LuckPillarNumber; got != nil {
		t.Errorf("2106 luck pillar = %d, want null", *got)
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

// The yearly table starts at the first supported birth year, which follows the
// calendar data.
func TestGetAnnualPillarsAtTheFirstSupportedYear(t *testing.T) {
	ict := time.FixedZone("ICT", 7*3600)
	chart, err := GetBaziChart(time.Date(1700, 6, 1, 12, 0, 0, 0, ict), ict, model.GenderMale)
	if err != nil {
		t.Fatal(err)
	}
	annual, err := GetAnnualPillars(chart, 1700, 3)
	if err != nil {
		t.Fatalf("1700 for 3 years: %v", err)
	}
	// 1700 is a Geng-Chen year and the birth falls after its Lunar New Year, so the nominal age starts at 1.
	if r := annual.AnnualPillars[0]; r.Year != 1700 || r.HeavenlyStem.Spelling != "geng" || r.EarthlyBranch.Spelling != "chen" || r.NominalAge != 1 {
		t.Errorf("1700 = %d %s/%s age %d", r.Year, r.HeavenlyStem.Spelling, r.EarthlyBranch.Spelling, r.NominalAge)
	}
	if _, err := GetAnnualPillars(chart, 1699, 3); !errors.Is(err, model.ErrInvalidYearRange) {
		t.Errorf("1699 precedes the first supported year: error = %v, want ErrInvalidYearRange", err)
	}
}
