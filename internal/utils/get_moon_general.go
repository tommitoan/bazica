package utils

import (
	"fmt"
	"time"

	"github.com/tommitoan/bazica/model"
)

// moonGeneralBranch is the Rat-first branch of the Moon General that starts at
// each principal term (zhongqi).
var moonGeneralBranch = map[string]int{
	model.MajorCold:      0,
	model.SpringShowers:  11,
	model.SpringEquinox:  10,
	model.GrainRain:      9,
	model.GrainBuds:      8,
	model.SummerSolstice: 7,
	model.MajorHeat:      6,
	model.EndOfHeat:      5,
	model.AutumnEquinox:  4,
	model.Frost:          3,
	model.MinorSnow:      2,
	model.WinterSolstice: 1,
}

// GetMoonGeneral returns the Rat-first branch index of the Moon General (Nguyet
// Tuong) for the calendar date of dateTime in its own location.
//
// The general changes at each principal term. The whole calendar day on which a
// term falls already belongs to the new general, so the Sun is read at the end
// of the day.
func GetMoonGeneral(dateTime time.Time) (int, error) {
	years, err := loadSolarTerms()
	if err != nil {
		return 0, err
	}
	y, m, d := dateTime.Date()
	endOfDay := time.Date(y, m, d, 23, 59, 59, 0, dateTime.Location())

	general, found := 0, false
	for _, year := range []int{endOfDay.Year() - 1, endOfDay.Year()} {
		terms, ok := years[year]
		if !ok {
			return 0, fmt.Errorf("%w: no solar term data for %d", model.ErrDateOutOfRange, year)
		}
		for _, term := range terms {
			branch, principal := moonGeneralBranch[term.name]
			if principal && !term.time.After(endOfDay) {
				general, found = branch, true
			}
		}
	}
	if !found {
		return 0, fmt.Errorf("%w: no principal term before %s", model.ErrDateOutOfRange, endOfDay)
	}
	return general, nil
}
