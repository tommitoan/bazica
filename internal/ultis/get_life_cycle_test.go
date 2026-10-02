package ultis

import (
	"testing"

	"github.com/tommitoan/bazica/model"
)

var lifeStages = []string{
	model.LC1, model.LC2, model.LC3, model.LC4, model.LC5, model.LC6,
	model.LC7, model.LC8, model.LC9, model.LC10, model.LC11, model.LC12,
}

// birthBranch maps each Heavenly Stem value to the branch value where its
// cycle starts at "Birth" (Chang Sheng).
var birthBranch = map[int]int{
	model.YangWoodValue:  model.PigValue,
	model.YinWoodValue:   model.HorseValue,
	model.YangFireValue:  model.TigerValue,
	model.YinFireValue:   model.RoosterValue,
	model.YangEarthValue: model.TigerValue,
	model.YinEarthValue:  model.RoosterValue,
	model.YangMetalValue: model.SnakeValue,
	model.YinMetalValue:  model.RatValue,
	model.YangWaterValue: model.MonkeyValue,
	model.YinWaterValue:  model.RabbitValue,
}

func mod12(n int) int { return ((n % 12) + 12) % 12 }

// Yang stems follow the branch order, Yin stems run against it.
func TestLifeCycleAllStemsAndBranches(t *testing.T) {
	for stem := 1; stem <= 10; stem++ {
		for branch := 1; branch <= 12; branch++ {
			step := branch - birthBranch[stem]
			if stem%2 == 0 {
				step = -step
			}
			want := lifeStages[mod12(step)]

			pillars := &model.FourPillars{
				YearPillar:  &model.YearPillar{EarthlyBranch: model.EarthlyBranch{Value: branch}},
				MonthPillar: &model.MonthPillar{EarthlyBranch: model.EarthlyBranch{Value: branch}},
				DayPillar:   &model.DayPillar{HeavenlyStem: model.HeavenlyStem{Value: stem}, EarthlyBranch: model.EarthlyBranch{Value: branch}},
				HourPillar:  &model.HourPillar{EarthlyBranch: model.EarthlyBranch{Value: branch}},
			}
			GetLifeCycleFromFourPillar(pillars)

			for name, got := range map[string]string{
				"year": pillars.YearPillar.LifeCycle, "month": pillars.MonthPillar.LifeCycle,
				"day": pillars.DayPillar.LifeCycle, "hour": pillars.HourPillar.LifeCycle,
			} {
				if got != want {
					t.Errorf("stem %d branch %d (%s pillar): got %q, want %q", stem, branch, name, got, want)
				}
			}
		}
	}
}

func TestLifeCycleKnownCases(t *testing.T) {
	tests := []struct {
		name         string
		stem, branch int
		want         string
	}{
		{"Yang Wood is born at Pig", model.YangWoodValue, model.PigValue, model.LC1},
		{"Yang Wood peaks at Rabbit", model.YangWoodValue, model.RabbitValue, model.LC5},
		{"Yin Wood is born at Horse", model.YinWoodValue, model.HorseValue, model.LC1},
		{"Yin Wood thrives at Rabbit", model.YinWoodValue, model.RabbitValue, model.LC4},
		{"Yin Wood peaks at Tiger", model.YinWoodValue, model.TigerValue, model.LC5},
		{"Yin Fire thrives at Horse", model.YinFireValue, model.HorseValue, model.LC4},
		{"Yin Metal is born at Rat", model.YinMetalValue, model.RatValue, model.LC1},
		{"Yin Water is born at Rabbit", model.YinWaterValue, model.RabbitValue, model.LC1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rule := GetLifeCycleRule(tc.stem)
			if got := GetLifeCycleByRuleAndBranch(tc.branch, rule, tc.stem%2 == 0); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
