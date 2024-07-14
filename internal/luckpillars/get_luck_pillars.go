package luckpillars

import (
	"fmt"
	"github.com/tommitoan/bazica/model"
	"log/slog"
)

func GetLuckPillars(fourPillars *model.FourPillars, gender int, prefixPath ...string) (*model.LuckPillars, error) {
	var luckPillars model.LuckPillars

	// Get 0 luck pillar
	luckPillars.LuckPillar0.HeavenlyStem = fourPillars.MonthPillar.HeavenlyStem
	luckPillars.LuckPillar0.EarthlyBranch = fourPillars.MonthPillar.EarthlyBranch
	luckPillars.LuckPillar0.Age = 0

	// Get Dayun rule (Increment or Decrement) with Yin/Yang Male/Female by Year Stem
	var incrementRule bool
	if (fourPillars.YearPillar.HeavenlyStem.Value+gender)%2 == 0 {
		incrementRule = true
	} else {
		incrementRule = false
	}
	slog.Info(fmt.Sprintf("LuckPillars increment rule: %t", incrementRule))

	// Calculate 13 luck pillars (without age and year)

	// Calculate age and year of luck pillar 1

	// Complete 11 other pillars

	return &luckPillars, nil
}
