package fourpillars

import (
	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
	"time"
)

// GetFourPillars computes the four pillars for the wall-clock time of dateTime
// in loc. A nil loc keeps the location dateTime already carries.
func GetFourPillars(dateTime time.Time, loc *time.Location) (*model.FourPillars, int, int, error) {
	if loc != nil {
		dateTime = dateTime.In(loc)
	}

	yearPillar, err := GetYearPillar(dateTime)
	if err != nil {
		return nil, 0, 0, err
	}
	yearPillar.GanZhi = utils.GanZhiOf(yearPillar.HeavenlyStem, yearPillar.EarthlyBranch)

	monthPillar, passed, remaining, err := GetMonthPillar(yearPillar, dateTime)
	if err != nil {
		return nil, 0, 0, err
	}
	monthPillar.GanZhi = utils.GanZhiOf(monthPillar.HeavenlyStem, monthPillar.EarthlyBranch)

	dayPillar, err := GetDayPillar(dateTime)
	if err != nil {
		return nil, 0, 0, err
	}
	dayPillar.GanZhi = utils.GanZhiOf(dayPillar.HeavenlyStem, dayPillar.EarthlyBranch)

	hourPillar, err := GetHourPillar(dayPillar, dateTime)
	if err != nil {
		return nil, 0, 0, err
	}
	hourPillar.GanZhi = utils.GanZhiOf(hourPillar.HeavenlyStem, hourPillar.EarthlyBranch)

	return &model.FourPillars{
		YearPillar:  yearPillar,
		MonthPillar: monthPillar,
		DayPillar:   dayPillar,
		HourPillar:  hourPillar,
	}, passed, remaining, nil
}
