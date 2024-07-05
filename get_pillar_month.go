package bazica

import (
	"github.com/tommitoan/bazica/model"
	"time"
)

func GetMonthPillar(path string, yearPillar *model.YearPillar, dateTime time.Time) (*model.MonthPillar, error) {
	var monthPillar model.MonthPillar
	monthPillar.Month = int(dateTime.Month())

	// Detect solar term
	termName, err := GetSolarTerm(path, dateTime)
	if err != nil {
		return nil, err
	}

	// Get earthly branch
	earthBranch := ConvertTermToBranch(termName)
	monthPillar.EarthlyBranch = earthBranch

	// Get heavenly stem
	valueOfFirstMonth := GetStemRuleByFireTigers(yearPillar.HeavenlyStem.Value)
	valueToCal := (valueOfFirstMonth - 1) + monthPillar.EarthlyBranch.Value
	if valueToCal > 10 {
		valueToCal = valueToCal - 10
	}
	heavenlyStem := CalculateHeavenlyStem(valueToCal)
	monthPillar.HeavenlyStem = heavenlyStem

	return &monthPillar, nil
}
