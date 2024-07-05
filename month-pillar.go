package bazica

import "time"

func GetMonthPillar(path string, yearPillar *YearPillar, dateTime time.Time) (*MonthPillar, error) {
	var monthPillar MonthPillar
	monthPillar.Month = int(dateTime.Month())

	// Detect solar term
	termName, err := DetectSolarTerm(path, dateTime)
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
