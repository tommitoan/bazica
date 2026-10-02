package fourpillars

import (
	"github.com/tommitoan/bazica/internal/utils"
	"github.com/tommitoan/bazica/model"
	"time"
)

func GetMonthPillar(yearPillar *model.YearPillar, dateTime time.Time) (*model.MonthPillar, int, int, error) {
	var monthPillar model.MonthPillar
	monthPillar.Month = int(dateTime.Month())

	// Detect solar term
	termName, passed, remaining, err := utils.GetSolarTerm(dateTime)
	if err != nil {
		return nil, 0, 0, err
	}

	// Get earthly branch
	earthBranch := utils.ConvertTermToBranch(termName)
	monthPillar.EarthlyBranch = earthBranch

	// Get heavenly stem
	valueOfFirstMonth := utils.GetStemRuleByFiveTigers(yearPillar.HeavenlyStem.Value)
	valueToCal := (valueOfFirstMonth - 1) + monthPillar.EarthlyBranch.Value
	if valueToCal > 10 {
		valueToCal = valueToCal - 10
	}
	heavenlyStem := utils.CalculateHeavenlyStem(valueToCal)
	monthPillar.HeavenlyStem = heavenlyStem

	return &monthPillar, passed, remaining, nil
}
