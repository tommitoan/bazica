package bazica

import (
	"github.com/tommitoan/bazica/v2/internal/analysis"
	"github.com/tommitoan/bazica/v2/internal/fourpillars"
	"github.com/tommitoan/bazica/v2/internal/luckpillars"
	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
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
// instant of Lichun (Start of Spring), and the month changes at the "initial" solar
// terms (jie). Analysis.YearReference carries display-only data, the Lichun time and
// the year of the lunar calendar; nothing is calculated from it.
// Dates inside the range returned by SupportedYears are supported; others return
// an error wrapping model.ErrDateOutOfRange.
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

	chart := &model.BaziChart{FourPillar: fourPillar, LuckPillars: luckPillars}
	if err := analysis.AttachAt(chart, gender, dateTime); err != nil {
		return nil, err
	}
	return chart, nil
}
