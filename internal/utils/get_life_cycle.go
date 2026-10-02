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
	offset := branchValue - 1
	if reverse {
		offset = -offset
	}
	valueToCal := (lifeCycleValue-1+offset)%12 + 1
	if valueToCal < 1 {
		valueToCal += 12
	}
	return ConvertValueToLifeCycle(valueToCal)
}
