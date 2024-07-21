package fourpillars

import (
	"encoding/json"
	"github.com/tommitoan/bazica/model"
	"log/slog"
	"time"
)

func GetFourPillars(dateTime time.Time, loc *time.Location, prefixPath ...string) (*model.FourPillars, int, int, error) {
	var path string
	if len(prefixPath) != 0 {
		path = prefixPath[0]
	}

	var fourPillars model.FourPillars

	// Get Year pillar
	yearPillar, err := GetYearPillar(path, dateTime)
	if err != nil {
		return nil, 0, 0, err
	}
	fourPillars.YearPillar = yearPillar

	// Get Month pillar
	monthPillar, passed, remaining, err := GetMonthPillar(path, yearPillar, dateTime)
	if err != nil {
		return nil, 0, 0, err
	}
	fourPillars.MonthPillar = monthPillar

	// Get Day pillar
	dayPillar, err := GetDayPillar(dateTime, loc)
	if err != nil {
		return nil, 0, 0, err
	}
	fourPillars.DayPillar = dayPillar

	// Get Hour pillar
	hourPillar, err := GetHourPillar(dayPillar, dateTime)
	if err != nil {
		return nil, 0, 0, err
	}
	fourPillars.HourPillar = hourPillar

	jsonData, _ := json.Marshal(fourPillars)
	slog.Info(string(jsonData))
	return &fourPillars, passed, remaining, nil
}
