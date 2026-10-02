package fourpillars

import (
	"github.com/tommitoan/bazica/internal/utils"
	"github.com/tommitoan/bazica/model"
	"time"
)

// dayCountEpoch is 1900-01-01, a Jia Xu day.
var dayCountEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

// GetDayPillar returns the day pillar for the wall-clock date of dateTime in its
// own location.
func GetDayPillar(dateTime time.Time) (*model.DayPillar, error) {
	var dayPillar model.DayPillar
	dayPillar.Day = dateTime.Day()

	// Count whole calendar days since the epoch. Using civil dates instead of
	// elapsed time keeps the count exact across offset changes.
	year, month, day := utils.BirthCalendarDate(dateTime)
	num := int(time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Sub(dayCountEpoch).Hours() / 24)

	// Get day stem
	stemRule := (num + 1) % 10
	if stemRule < 1 {
		stemRule = stemRule + 10
	}
	dayPillar.HeavenlyStem = utils.CalculateHeavenlyStem(stemRule)

	// Get day branch
	branchRule := (num - 3) % 12
	if branchRule < 1 {
		branchRule = branchRule + 12
	}
	dayPillar.EarthlyBranch = utils.CalculateEarthlyBranch(branchRule)

	return &dayPillar, nil
}
