package fourpillars

import (
	"github.com/tommitoan/bazica/internal/ultis"
	"github.com/tommitoan/bazica/model"
	"time"
)

func GetMonthPillar(path string, yearPillar *model.YearPillar, dateTime time.Time) (*model.MonthPillar, error) {
	var monthPillar model.MonthPillar
	monthPillar.Month = int(dateTime.Month())

	// Detect solar term
	termName, err := ultis.GetSolarTerm(path, dateTime)
	if err != nil {
		return nil, err
	}

	// Get earthly branch
	earthBranch := ultis.ConvertTermToBranch(termName)
	monthPillar.EarthlyBranch = earthBranch

	// Get heavenly stem
	valueOfFirstMonth := ultis.GetStemRuleByFireTigers(yearPillar.HeavenlyStem.Value)
	valueToCal := (valueOfFirstMonth - 1) + monthPillar.EarthlyBranch.Value
	if valueToCal > 10 {
		valueToCal = valueToCal - 10
	}
	heavenlyStem := ultis.CalculateHeavenlyStem(valueToCal)
	monthPillar.HeavenlyStem = heavenlyStem

	return &monthPillar, nil
}
