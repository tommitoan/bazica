package bazica

import (
	"github.com/tommitoan/bazica/internal/fourpillars"
	"github.com/tommitoan/bazica/internal/luckpillars"
	"github.com/tommitoan/bazica/internal/utils"
	"github.com/tommitoan/bazica/model"
	"time"
)

// GetBaziChart calculates the Ba-zi chart (four pillars and luck pillars) for a
// birth time.
//
// dateTime is read as wall-clock time in loc; pass nil to use the location
// dateTime already carries. gender is model.GenderFemale (0) or model.GenderMale
// (1) and decides the direction of the luck pillars.
//
// Conventions: the day changes at 23:00 (the Rat hour), the year changes at the
// Lunar New Year, and the month changes at the "initial" solar terms (jie).
// Dates from 1900-01-01 to 2099-12-31 are supported; others return an error
// wrapping model.ErrDateOutOfRange.
func GetBaziChart(dateTime time.Time, loc *time.Location, gender int) (*model.BaziChart, error) {
	if gender != model.GenderFemale && gender != model.GenderMale {
		return nil, model.ErrInvalidGender
	}
	if loc != nil {
		dateTime = dateTime.In(loc)
	}

	fourPillar, passed, remaining, err := fourpillars.GetFourPillars(dateTime, loc)
	if err != nil {
		return nil, err
	}
	utils.GetLifeCycleFromFourPillar(fourPillar)

	luckPillars, err := luckpillars.GetLuckPillars(fourPillar, gender, passed, remaining, dateTime)
	if err != nil {
		return nil, err
	}

	return &model.BaziChart{FourPillar: fourPillar, LuckPillars: luckPillars}, nil
}
