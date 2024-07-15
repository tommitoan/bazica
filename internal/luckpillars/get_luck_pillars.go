package luckpillars

import (
	"fmt"
	"github.com/tommitoan/bazica/internal/ultis"
	"github.com/tommitoan/bazica/model"
	"log/slog"
)

func GetLuckPillars(fourPillars *model.FourPillars, gender int, prefixPath ...string) (*model.LuckPillars, error) {
	luckPillars := model.LuckPillars{}

	// Get 0 luck pillar
	zeroLuckPillars := &model.LuckPillar{
		HeavenlyStem:  fourPillars.MonthPillar.HeavenlyStem,
		EarthlyBranch: fourPillars.MonthPillar.EarthlyBranch,
		Age:           0,
		Number:        0,
	}

	luckPillars.LuckPillars = append(luckPillars.LuckPillars, zeroLuckPillars)

	// Get Dayun rule (Increment or Decrement) with Yin/Yang Male/Female by Year Stem
	var incrementRule int
	if (fourPillars.YearPillar.HeavenlyStem.Value+gender)%2 == 0 {
		// Yang Male or Yin Female
		incrementRule = 1
	} else {
		// Yin Male or Yang Female
		incrementRule = -1
	}
	slog.Info(fmt.Sprintf("LuckPillars increment rule: %v", incrementRule))

	// Calculate the remaining 11 luck pillars (excluding age and year)
	var branchValue int = fourPillars.MonthPillar.EarthlyBranch.Value
	var stemValue int = fourPillars.MonthPillar.HeavenlyStem.Value
	// pillar1
	for i := 1; i < 12; i++ {
		branchValue += incrementRule
		if branchValue > 12 {
			branchValue = branchValue - 12
		} else if branchValue < 1 {
			branchValue = branchValue + 12
		}
		stemValue += incrementRule
		if stemValue > 10 {
			stemValue = stemValue - 10
		} else if stemValue < 1 {
			stemValue = stemValue + 10
		}

		branch := ultis.CalculateEarthlyBranch(branchValue)
		stem := ultis.CalculateHeavenlyStem(stemValue)

		tempLuckPillars := &model.LuckPillar{
			HeavenlyStem:  stem,
			EarthlyBranch: branch,
			Number:        i,
		}

		luckPillars.LuckPillars = append(luckPillars.LuckPillars, tempLuckPillars)
	}

	// Calculate age and year of luck pillar 1

	// Complete 11 other pillars

	return &luckPillars, nil
}
