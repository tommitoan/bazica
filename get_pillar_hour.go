package bazica

import (
	"github.com/tommitoan/bazica/model"
	"math"
	"time"
)

func GetHourPillar(dayPillar *model.DayPillar, dateTime time.Time) (*model.HourPillar, error) {
	var hourPillar model.HourPillar
	hourPillar.Hour = model.TimeOfDay{
		Hour:   dateTime.Hour(),
		Minute: dateTime.Minute(),
	}

	// Get earthly branch
	dateHour := dateTime.Hour()
	branch := CalculateEarthlyBranch(roundInt(dateHour / 2))
	hourPillar.EarthlyBranch = branch

	// Get heavenly stem
	valueOfRatHour := GetStemRuleByFireRats(dayPillar.HeavenlyStem.Value)
	var stem int
	diff := hourPillar.EarthlyBranch.Value - 11
	if diff < 0 {
		diff = diff + 12
	}
	stem = valueOfRatHour + diff
	if stem > 10 {
		stem = stem - 10
	}

	heavenlyStem := CalculateHeavenlyStem(stem)
	hourPillar.HeavenlyStem = heavenlyStem

	return &hourPillar, nil
}

func roundInt(x int) int {
	return int(math.Round(float64(x)))
}