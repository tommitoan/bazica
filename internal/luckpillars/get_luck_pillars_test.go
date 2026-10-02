package luckpillars

import (
	"testing"
	"time"

	"github.com/tommitoan/bazica/internal/ultis"
	"github.com/tommitoan/bazica/model"
)

const (
	female = 0
	male   = 1
)

// chart builds the minimal four pillars the luck pillar calculation reads.
func chart(yearStem, monthStem, monthBranch int) *model.FourPillars {
	return &model.FourPillars{
		YearPillar: &model.YearPillar{HeavenlyStem: ultis.CalculateHeavenlyStem(yearStem)},
		MonthPillar: &model.MonthPillar{
			HeavenlyStem:  ultis.CalculateHeavenlyStem(monthStem),
			EarthlyBranch: ultis.CalculateEarthlyBranch(monthBranch),
		},
	}
}

func TestGetLuckPillarsStartTime(t *testing.T) {
	// Yin Wood year, male: the count runs backwards from the previous term,
	// 4124 minutes = 68.7 hours = 343 days.
	fourPillars := chart(model.YinWoodValue, model.YangWaterValue, model.HorseValue)
	birth := time.Date(1995, time.June, 8, 22, 5, 0, 0, time.UTC)

	got, err := GetLuckPillars(fourPillars, male, 4124, 21429, birth)
	if err != nil {
		t.Fatalf("GetLuckPillars() error = %v", err)
	}
	want := time.Date(1996, time.May, 18, 22, 5, 0, 0, time.UTC)
	if !got.LuckPillars[1].Time.Equal(want) {
		t.Errorf("first luck pillar starts %v, want %v", got.LuckPillars[1].Time, want)
	}
}

// Yang year + male or Yin year + female run forward using the minutes until
// the next term; the other two combinations run backward using the minutes
// since the previous term.
func TestGetLuckPillarsDirection(t *testing.T) {
	const (
		oneYear  = 3 * 24 * 60 // 3 days of minutes map to one year of age
		twoYears = 2 * oneYear
	)
	birth := time.Date(2000, time.March, 10, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		yearStem     int
		gender       int
		wantForward  bool
		wantStartAge int
	}{
		{"Yang year, male", model.YangWoodValue, male, true, 2},
		{"Yang year, female", model.YangWoodValue, female, false, 1},
		{"Yin year, male", model.YinWoodValue, male, false, 1},
		{"Yin year, female", model.YinWoodValue, female, true, 2},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Month pillar Ren Wu (stem 9, branch 5).
			got, err := GetLuckPillars(chart(tc.yearStem, model.YangWaterValue, model.HorseValue), tc.gender, oneYear, twoYears, birth)
			if err != nil {
				t.Fatal(err)
			}
			first, second := got.LuckPillars[1], got.LuckPillars[2]

			wantStem, wantBranch := model.YinWaterValue, model.GoatValue // Gui Wei
			wantStem2, wantBranch2 := model.YangWoodValue, model.MonkeyValue
			if !tc.wantForward {
				wantStem, wantBranch = model.YinMetalValue, model.SnakeValue // Xin Si
				wantStem2, wantBranch2 = model.YangMetalValue, model.DragonValue
			}
			if first.HeavenlyStem.Value != wantStem || first.EarthlyBranch.Value != wantBranch {
				t.Errorf("pillar 1 = %s %s, want stem %d branch %d", first.HeavenlyStem.Name, first.EarthlyBranch.Name, wantStem, wantBranch)
			}
			if second.HeavenlyStem.Value != wantStem2 || second.EarthlyBranch.Value != wantBranch2 {
				t.Errorf("pillar 2 = %s %s, want stem %d branch %d", second.HeavenlyStem.Name, second.EarthlyBranch.Name, wantStem2, wantBranch2)
			}
			if want := birth.AddDate(tc.wantStartAge, 0, 0); !first.Time.Equal(want) {
				t.Errorf("starts %v, want %v", first.Time, want)
			}
		})
	}
}

func TestGetLuckPillarsWrapAround(t *testing.T) {
	birth := time.Date(2000, time.January, 20, 12, 0, 0, 0, time.UTC)

	// Forward from Gui Chou wraps to Jia Yin.
	got, err := GetLuckPillars(chart(model.YangWoodValue, model.YinWaterValue, model.OxValue), male, 100, 100, birth)
	if err != nil {
		t.Fatal(err)
	}
	if p := got.LuckPillars[1]; p.HeavenlyStem.Value != model.YangWoodValue || p.EarthlyBranch.Value != model.TigerValue {
		t.Errorf("forward wrap = %d/%d, want 1/1", p.HeavenlyStem.Value, p.EarthlyBranch.Value)
	}

	// Backward from Jia Yin wraps to Gui Chou.
	got, err = GetLuckPillars(chart(model.YangWoodValue, model.YangWoodValue, model.TigerValue), female, 100, 100, birth)
	if err != nil {
		t.Fatal(err)
	}
	if p := got.LuckPillars[1]; p.HeavenlyStem.Value != model.YinWaterValue || p.EarthlyBranch.Value != model.OxValue {
		t.Errorf("backward wrap = %d/%d, want 10/12", p.HeavenlyStem.Value, p.EarthlyBranch.Value)
	}
}

func TestGetLuckPillarsPeriods(t *testing.T) {
	birth := time.Date(1990, time.August, 15, 8, 0, 0, 0, time.UTC)
	got, err := GetLuckPillars(chart(model.YangMetalValue, model.YangFireValue, model.MonkeyValue), male, 5000, 30000, birth)
	if err != nil {
		t.Fatal(err)
	}
	pillars := got.LuckPillars
	if len(pillars) != 12 {
		t.Fatalf("got %d pillars, want 12", len(pillars))
	}
	for i, p := range pillars {
		if p.Number != i {
			t.Errorf("pillar %d has number %d", i, p.Number)
		}
		if p.GanZhi.Name == "" {
			t.Errorf("pillar %d has no GanZhi", i)
		}
		if i == 0 {
			continue
		}
		if i >= 2 {
			if p.YearStart != pillars[i-1].YearStart+10 {
				t.Errorf("pillar %d starts %d, want %d", i, p.YearStart, pillars[i-1].YearStart+10)
			}
		}
		if p.YearEnd != p.YearStart+9 {
			t.Errorf("pillar %d spans %d-%d, want 10 years", i, p.YearStart, p.YearEnd)
		}
		if i >= 1 && pillars[i-1].YearEnd >= p.YearStart {
			t.Errorf("pillar %d overlaps pillar %d (%d >= %d)", i, i-1, pillars[i-1].YearEnd, p.YearStart)
		}
	}
	if pillars[0].YearStart != birth.Year() {
		t.Errorf("pillar 0 starts %d, want %d", pillars[0].YearStart, birth.Year())
	}
}
