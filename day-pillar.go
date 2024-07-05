package bazica

import (
	"time"
)

func GetDayPillar(dateTime time.Time, loc *time.Location) (*DayPillar, error) {
	var dayPillar DayPillar
	dayPillar.Day = dateTime.Day()

	// From 23:00 is new day (Rat hour)
	dateTime = dateTime.Add(time.Hour)
	
	milestone := time.Date(1900, 1, 1, 0, 0, 0, 0, loc)
	num := dateTime.Sub(milestone).Hours() / 24.0

	// Get day stem
	stemRule := (int(num) + 1) % 10
	stem := CalculateHeavenlyStem(stemRule)
	dayPillar.HeavenlyStem = stem

	// Get day branch
	branchRule := (int(num) - 3) % 12
	branch := CalculateEarthlyBranch(branchRule)
	dayPillar.EarthlyBranch = branch

	return &dayPillar, nil
}