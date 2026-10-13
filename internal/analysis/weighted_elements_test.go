package analysis

import (
	"math"
	"testing"

	"github.com/tommitoan/bazica/v2/model"
)

func TestWeightedElementsThreeHorses(t *testing.T) {
	chart := chartOf(chart1995)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	w := chart.Analysis.WeightedElements
	if w.Scheme != "hidden_stems_v1" || w.Total != 8 || len(w.Elements) != 5 {
		t.Fatalf("weighted metadata = %+v", w)
	}
	fire := w.Elements["fire"]
	if fire.Score != 2.1 || fire.Level != 2 || fire.DayMaster != 0 || len(fire.Contributions) != 3 {
		t.Errorf("three horses fire = %+v, want 2.1 direct officer, level 2", fire)
	}
	if len(fire.TenGods) != 2 || fire.TenGods[0].TenGod.Code != "ten_god.direct_officer" || fire.TenGods[0].Score != 2.1 ||
		fire.TenGods[1].TenGod.Code != "ten_god.seven_killings" || fire.TenGods[1].Score != 0 {
		t.Errorf("fire Ten Gods = %+v", fire.TenGods)
	}
	for _, c := range fire.Contributions {
		if c.Source != "hidden_stem" || c.EarthlyBranch == nil || c.EarthlyBranch.Value != 5 || c.Score != 0.7 {
			t.Errorf("horse contribution = %+v", c)
		}
	}
	earth := w.Elements["earth"]
	if earth.Score != 0.9 || earth.TenGods[0].TenGod.Code != "ten_god.direct_resource" || earth.TenGods[0].Score != 0.9 {
		t.Errorf("three horses earth = %+v", earth)
	}
	metal := w.Elements["metal"]
	if metal.Score != 1 || metal.DayMaster != 1 || len(metal.TenGods) != 2 || metal.TenGods[0].Score != 0 || metal.TenGods[1].Score != 0 {
		t.Errorf("day master accounting = %+v", metal)
	}
	var total, gods, dayMaster float64
	for _, element := range w.Elements {
		total += element.Score
		dayMaster += element.DayMaster
		var parts float64
		for _, c := range element.Contributions {
			parts += c.Score
		}
		if math.Abs(parts-element.Score) > 1e-9 {
			t.Errorf("contributions %.2f != score %.2f", parts, element.Score)
		}
		for _, god := range element.TenGods {
			gods += god.Score
		}
	}
	if math.Abs(total-8) > 1e-9 || math.Abs(gods-7) > 1e-9 || dayMaster != 1 {
		t.Errorf("totals: elements %.2f, Ten Gods %.2f, Day Master %.2f", total, gods, dayMaster)
	}
}

func TestWeightedLevels(t *testing.T) {
	for _, tc := range []struct{ tenths, level int }{{0, 0}, {4, 0}, {5, 1}, {14, 1}, {15, 2}, {24, 2}, {25, 3}, {34, 3}, {35, 4}, {80, 4}} {
		if got := weightedLevel(tc.tenths); got != tc.level {
			t.Errorf("%d tenths: level %d, want %d", tc.tenths, got, tc.level)
		}
	}
}
