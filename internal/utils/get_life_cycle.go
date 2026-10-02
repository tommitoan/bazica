package utils

import "github.com/tommitoan/bazica/model"

func GetLifeCycleFromFourPillar(pillars *model.FourPillars) *model.FourPillars {
	dayValue := pillars.DayPillar.HeavenlyStem.Value
	lifeCycleValue := GetLifeCycleRule(dayValue)
	// Yang stems advance through the twelve stages along the branches while
	// Yin stems move against them.
	reverse := dayValue%2 == 0

	pillars.YearPillar.LifeCycle = GetLifeCycleByRuleAndBranch(pillars.YearPillar.EarthlyBranch.Value, lifeCycleValue, reverse)
	pillars.MonthPillar.LifeCycle = GetLifeCycleByRuleAndBranch(pillars.MonthPillar.EarthlyBranch.Value, lifeCycleValue, reverse)
	pillars.DayPillar.LifeCycle = GetLifeCycleByRuleAndBranch(pillars.DayPillar.EarthlyBranch.Value, lifeCycleValue, reverse)
	pillars.HourPillar.LifeCycle = GetLifeCycleByRuleAndBranch(pillars.HourPillar.EarthlyBranch.Value, lifeCycleValue, reverse)

	return pillars
}

// GetLifeCycleByRuleAndBranch returns the life stage of a branch, given the
// stage of the Tiger branch (lifeCycleValue) and the direction of the cycle.
func GetLifeCycleByRuleAndBranch(branchValue int, lifeCycleValue int, reverse bool) string {
	return ConvertValueToLifeCycle(GetLifeCycleValueByRuleAndBranch(branchValue, lifeCycleValue, reverse))
}

// GetLifeCycleValueByRuleAndBranch is GetLifeCycleByRuleAndBranch returning the
// stage value (1-12) instead of its name.
func GetLifeCycleValueByRuleAndBranch(branchValue int, lifeCycleValue int, reverse bool) int {
	offset := branchValue - 1
	if reverse {
		offset = -offset
	}
	valueToCal := (lifeCycleValue-1+offset)%12 + 1
	if valueToCal < 1 {
		valueToCal += 12
	}
	return valueToCal
}

// GetLifeCycleValueOfStem returns the life stage value (1-12) of any Heavenly
// Stem at a branch. Yang stems advance along the branches, Yin stems run against them.
func GetLifeCycleValueOfStem(stemValue int, branchValue int) int {
	return GetLifeCycleValueByRuleAndBranch(branchValue, GetLifeCycleRule(stemValue), stemValue%2 == 0)
}
