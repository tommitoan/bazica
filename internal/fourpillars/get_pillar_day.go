package fourpillars

import (
	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
	"time"
)

// dayCountEpoch is 1900-01-01, a Jia Xu day.
var dayCountEpoch = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

// daysSinceEpoch counts whole calendar days from the epoch to a civil date. It
// subtracts Unix seconds of two UTC midnights, which is exact for any year;
// time.Time.Sub would saturate at about 292 years and give a constant day pillar
// outside 1607-2192. Using civil dates instead of elapsed time also keeps the
// count exact across offset changes.
func daysSinceEpoch(year int, month time.Month, day int) int {
	midnight := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return int((midnight.Unix() - dayCountEpoch.Unix()) / secondsPerDay)
}

const secondsPerDay = 24 * 60 * 60

// GetDayPillar returns the day pillar for the wall-clock date of dateTime in its
// own location.
func GetDayPillar(dateTime time.Time) (*model.DayPillar, error) {
	var dayPillar model.DayPillar
	dayPillar.Day = dateTime.Day()

	year, month, day := utils.BirthCalendarDate(dateTime)
	num := daysSinceEpoch(year, month, day)

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
