package bazica

import (
	"fmt"
)

func GetFourPillarChart(path, input string) (any, error) {
	// Get Year pillar
	yearPillar, err := GetYearPillar(path, input)
	if err != nil {
		return nil, err
	}

	// Get Month pillar
	monthPillar, err := GetMonthPillar(path, input, yearPillar)
	if err != nil {
		return nil, err
	}

	fmt.Println(monthPillar)

	return nil, nil
}
