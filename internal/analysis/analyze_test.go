package analysis

import (
	"errors"
	"testing"

	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

// chartOf builds a chart from (stem, branch) values in year, month, day, hour order.
func chartOf(pairs [4][2]int) *model.BaziChart {
	mk := func(i int) (model.HeavenlyStem, model.EarthlyBranch, model.GanZhi) {
		s, b := utils.CalculateHeavenlyStem(pairs[i][0]), utils.CalculateEarthlyBranch(pairs[i][1])
		return s, b, utils.GanZhiOf(s, b)
	}
	ys, yb, yg := mk(0)
	ms, mb, mg := mk(1)
	ds, db, dg := mk(2)
	hs, hb, hg := mk(3)
	return &model.BaziChart{FourPillar: &model.FourPillars{
		YearPillar:  &model.YearPillar{HeavenlyStem: ys, EarthlyBranch: yb, GanZhi: yg},
		MonthPillar: &model.MonthPillar{HeavenlyStem: ms, EarthlyBranch: mb, GanZhi: mg},
		DayPillar:   &model.DayPillar{HeavenlyStem: ds, EarthlyBranch: db, GanZhi: dg},
		HourPillar:  &model.HourPillar{HeavenlyStem: hs, EarthlyBranch: hb, GanZhi: hg},
	}}
}

// The 1995-06-08 reference chart: Yi-Hai, Ren-Wu, Geng-Wu, Ren-Wu.
var chart1995 = [4][2]int{{2, 10}, {9, 5}, {7, 5}, {9, 5}}

func TestAttachFillsEveryPillarAndTheChart(t *testing.T) {
	chart := chartOf(chart1995)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatalf("Attach() error = %v", err)
	}
	fp := chart.FourPillar
	for name, a := range map[string]*model.PillarAnalysis{
		"year": fp.YearPillar.Analysis, "month": fp.MonthPillar.Analysis,
		"day": fp.DayPillar.Analysis, "hour": fp.HourPillar.Analysis,
	} {
		if a == nil {
			t.Fatalf("%s pillar has no analysis", name)
		}
		if a.Stars == nil {
			t.Errorf("%s stars are nil, want a list", name)
		}
		if len(a.HiddenStems) == 0 {
			t.Errorf("%s pillar has no hidden stems", name)
		}
	}
	if chart.Analysis == nil {
		t.Fatal("chart has no analysis")
	}
}

func TestAttachDayMasterPillar(t *testing.T) {
	chart := chartOf(chart1995)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	day := chart.FourPillar.DayPillar.Analysis
	if !day.IsDayMaster || day.TenGod != nil || day.IsVoid {
		t.Errorf("day pillar analysis = %+v, want Day Master, no Ten God, not void", day)
	}
	year := chart.FourPillar.YearPillar.Analysis
	if year.IsDayMaster || year.TenGod == nil || year.TenGod.Code != "ten_god.direct_wealth" {
		t.Errorf("year pillar Ten God = %+v, want direct wealth", year.TenGod)
	}
	// Day Master stage must agree with the existing LifeCycle names.
	if year.DayMasterStage.EN != "Sick" {
		t.Errorf("year Day Master stage = %q, want Sick", year.DayMasterStage.EN)
	}
}

func TestAttachVoidAndClashFlags(t *testing.T) {
	chart := chartOf(chart1995)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	fp := chart.FourPillar
	// Void branches of Geng-Wu are Xu and Hai: only the Hai year pillar is void.
	if !fp.YearPillar.Analysis.IsVoid || fp.MonthPillar.Analysis.IsVoid || fp.HourPillar.Analysis.IsVoid {
		t.Error("only the year pillar (Hai) should be void")
	}
	for name, a := range map[string]*model.PillarAnalysis{
		"year": fp.YearPillar.Analysis, "month": fp.MonthPillar.Analysis,
		"day": fp.DayPillar.Analysis, "hour": fp.HourPillar.Analysis,
	} {
		if a.HeavenEarthClash {
			t.Errorf("%s pillar flagged, but the 1995 chart has no clash", name)
		}
	}
}

func TestAttachFlagsAClashingChart(t *testing.T) {
	// Jia-Zi year, Geng-Wu month, Geng-Zi day, Bing-Zi hour.
	// Geng overcomes Jia and Wu opposes Zi, so the month is flagged; Bing overcomes
	// Geng and Zi opposes Wu, so the hour is flagged. The overcome cells (year and
	// day) are not flagged by those pairs.
	chart := chartOf([4][2]int{{1, 11}, {7, 5}, {7, 11}, {3, 11}})
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	fp := chart.FourPillar
	for name, tc := range map[string]struct {
		got  bool
		want bool
	}{
		"year":  {fp.YearPillar.Analysis.HeavenEarthClash, false},
		"month": {fp.MonthPillar.Analysis.HeavenEarthClash, true},
		"day":   {fp.DayPillar.Analysis.HeavenEarthClash, false},
		"hour":  {fp.HourPillar.Analysis.HeavenEarthClash, true},
	} {
		if tc.got != tc.want {
			t.Errorf("%s pillar clash = %v, want %v", name, tc.got, tc.want)
		}
	}
}

func TestAttachChartLevelFacts(t *testing.T) {
	chart := chartOf(chart1995)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	a := chart.Analysis
	if a.DayMaster.HeavenlyStem.Value != 7 || a.DayMaster.Element.Code != "element.metal" || a.DayMaster.Polarity.Code != "polarity.yang" {
		t.Errorf("day master = %+v", a.DayMaster)
	}
	if len(a.VoidBranches) != 2 || a.VoidBranches[0].Value != 9 || a.VoidBranches[1].Value != 10 {
		t.Errorf("void branches = %+v", a.VoidBranches)
	}
	for name, p := range map[string]model.AuxiliaryPillar{"thai nguyen": a.ThaiNguyen, "thai tuc": a.ThaiTuc, "life palace": a.LifePalace} {
		if p.Nayin.Code == "" || p.GanZhi.Name != p.Nayin.EN {
			t.Errorf("%s: Nayin %+v does not match GanZhi %q", name, p.Nayin, p.GanZhi.Name)
		}
	}
	if a.LifePalace.HeavenlyStem.Value != 8 || a.LifePalace.EarthlyBranch.Value != 4 {
		t.Errorf("life palace = %d/%d, want Xin-Si (8/4)", a.LifePalace.HeavenlyStem.Value, a.LifePalace.EarthlyBranch.Value)
	}
	if a.DayMasterStrength != nil || a.UsefulGod != nil {
		t.Error("reserved fields must stay nil")
	}
}

func TestAttachElementCounts(t *testing.T) {
	chart := chartOf(chart1995)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	c := chart.Analysis.ElementCounts
	// Stems Yi Ren Geng Ren: wood 1, metal 1, water 2. Branches Hai Wu Wu Wu: water 1, fire 3.
	if c.Stems != (model.ElementCount{Wood: 1, Metal: 1, Water: 2}) {
		t.Errorf("stems = %+v", c.Stems)
	}
	if c.Branches != (model.ElementCount{Fire: 3, Water: 1}) {
		t.Errorf("branches = %+v", c.Branches)
	}
	total := c.HiddenStems.Wood + c.HiddenStems.Fire + c.HiddenStems.Earth + c.HiddenStems.Metal + c.HiddenStems.Water
	if total != 2+2+2+2 {
		t.Errorf("hidden stems total = %d, want 8 (Hai 2, Wu 2 x3)", total)
	}
}

func TestAttachRejectsInvalidInput(t *testing.T) {
	cases := map[string]*model.BaziChart{
		"no pillars": {},
		"missing hour pillar": func() *model.BaziChart {
			c := chartOf(chart1995)
			c.FourPillar.HourPillar = nil
			return c
		}(),
		"stem out of range": func() *model.BaziChart {
			c := chartOf(chart1995)
			c.FourPillar.DayPillar.HeavenlyStem.Value = 11
			return c
		}(),
		"branch out of range": func() *model.BaziChart {
			c := chartOf(chart1995)
			c.FourPillar.MonthPillar.EarthlyBranch.Value = 0
			return c
		}(),
		"unknown Nayin": func() *model.BaziChart {
			c := chartOf(chart1995)
			c.FourPillar.YearPillar.GanZhi.Name = "Unknown"
			return c
		}(),
	}
	for name, chart := range cases {
		if err := Attach(chart, model.GenderMale); !errors.Is(err, errInvalidPillar) {
			t.Errorf("%s: error = %v, want errInvalidPillar", name, err)
		}
	}
}
