package analysis

import (
	"errors"
	"testing"

	"github.com/tommitoan/bazica/model"
)

func TestCycleStemBranchKnownYears(t *testing.T) {
	// 1984 Jia-Zi, 1995 Yi-Hai, 2026 Bing-Wu, 1900 Geng-Zi, 2099 Ji-Wei.
	for _, tc := range []struct{ year, stem, branch int }{
		{1984, 1, 11}, {1995, 2, 10}, {2026, 3, 5}, {1900, 7, 11}, {2099, 6, 6},
	} {
		if s, b := cycleStemBranch(tc.year); s != tc.stem || b != tc.branch {
			t.Errorf("cycleStemBranch(%d) = %d/%d, want %d/%d", tc.year, s, b, tc.stem, tc.branch)
		}
	}
}

func TestCycleRepeatsEverySixtyYears(t *testing.T) {
	for year := 1900; year+60 <= 2099; year++ {
		s1, b1 := cycleStemBranch(year)
		s2, b2 := cycleStemBranch(year + 60)
		if s1 != s2 || b1 != b2 {
			t.Fatalf("year %d and %d differ", year, year+60)
		}
	}
}

func TestYearPillarYear(t *testing.T) {
	// The model's Year is the civil birth year; the pillar's own year is one less
	// for a birth before the Lunar New Year.
	pillar := func(civil, pillarYear int) *model.YearPillar {
		s, b := cycleStemBranch(pillarYear)
		return &model.YearPillar{Year: civil, HeavenlyStem: model.HeavenlyStem{Value: s}, EarthlyBranch: model.EarthlyBranch{Value: b}}
	}
	for _, tc := range []struct {
		name          string
		civil, pillar int
	}{
		{"birth after Lunar New Year", 1995, 1995},
		{"January birth before Lunar New Year", 2000, 1999},
		{"first covered year", 1900, 1900},
	} {
		if got := yearPillarYear(pillar(tc.civil, tc.pillar)); got != tc.pillar {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.pillar)
		}
	}
}

func TestAnnualPillarsRequiresAChart(t *testing.T) {
	if _, err := AnnualPillars(nil, 1995, 1); !errors.Is(err, model.ErrInvalidYearRange) {
		t.Errorf("error = %v, want ErrInvalidYearRange", err)
	}
	if _, err := AnnualPillars(&model.BaziChart{}, 1995, 1); !errors.Is(err, errInvalidPillar) {
		t.Errorf("error = %v, want errInvalidPillar", err)
	}
}

func TestAttachLuckSkipsAChartWithoutLuckPillars(t *testing.T) {
	chart := chartOf(chart1995)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatalf("Attach() without luck pillars: %v", err)
	}
}

func TestAttachLuckRejectsInvalidLuckPillars(t *testing.T) {
	for name, mutate := range map[string]func(*model.LuckPillar){
		"stem out of range": func(lp *model.LuckPillar) { lp.HeavenlyStem.Value = 0 },
		"unknown Nayin":     func(lp *model.LuckPillar) { lp.GanZhi.Name = "Unknown" },
	} {
		chart := chartOf(chart1995)
		chart.FourPillar.YearPillar.Year = 1995
		lp := &model.LuckPillar{Number: 1, HeavenlyStem: model.HeavenlyStem{Value: 8}, EarthlyBranch: model.EarthlyBranch{Value: 4}, YearStart: 1996, YearEnd: 2005}
		lp.GanZhi.Name = "Wax metal"
		mutate(lp)
		chart.LuckPillars = &model.LuckPillars{LuckPillars: []*model.LuckPillar{lp}}
		if err := Attach(chart, model.GenderMale); !errors.Is(err, errInvalidPillar) {
			t.Errorf("%s: error = %v, want errInvalidPillar", name, err)
		}
	}
}

func TestAttachLuckAgesAndFlags(t *testing.T) {
	chart := chartOf(chart1995)
	chart.FourPillar.YearPillar.Year = 1995
	// Xin-Si luck pillar from 1996: nominal age 2 to 11. Xin overcomes Yi and Si opposes Hai
	// (the year pillar), so it clashes with the natal chart.
	lp := &model.LuckPillar{Number: 1, HeavenlyStem: model.HeavenlyStem{Value: 8}, EarthlyBranch: model.EarthlyBranch{Value: 4}, YearStart: 1996, YearEnd: 2005}
	lp.GanZhi.Name = "Wax metal"
	chart.LuckPillars = &model.LuckPillars{LuckPillars: []*model.LuckPillar{lp}}
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	a := lp.Analysis
	if a == nil || a.AgeStart != 2 || a.AgeEnd != 11 {
		t.Fatalf("analysis = %+v, want ages 2 to 11", a)
	}
	if !a.HeavenEarthClash {
		t.Error("Xin-Si should clash with the Yi-Hai year pillar")
	}
	if a.TenGod.Code != "ten_god.rob_wealth" || a.StemStageAtOwnBranch.Code != "stage.death" {
		t.Errorf("ten god %s, stage %s", a.TenGod.Code, a.StemStageAtOwnBranch.Code)
	}
	if len(a.HiddenStems) != 3 || a.HiddenStems[0].StageAtOwnBranch.Code != "stage.thriving" {
		t.Errorf("hidden stems = %+v", a.HiddenStems)
	}
}
