package fourpillars

import (
	"encoding/json"
	"github.com/tommitoan/bazica/model"
	"log/slog"
	"time"
)

func GetFourPillars(dateTime time.Time, loc *time.Location, prefixPath ...string) (*model.FourPillars, error) {
	var path string
	if len(prefixPath) != 0 {
		path = prefixPath[0]
	}

	var baziChart model.FourPillars

	// Get Year pillar
	yearPillar, err := GetYearPillar(path, dateTime)
	if err != nil {
		return nil, err
	}
	baziChart.YearPillar = yearPillar

	// Get Month pillar
	monthPillar, err := GetMonthPillar(path, yearPillar, dateTime)
	if err != nil {
		return nil, err
	}
	baziChart.MonthPillar = monthPillar

	// Get Day pillar
	dayPillar, err := GetDayPillar(dateTime, loc)
	if err != nil {
		return nil, err
	}
	baziChart.DayPillar = dayPillar

	// Get Hour pillar
	hourPillar, err := GetHourPillar(dayPillar, dateTime)
	if err != nil {
		return nil, err
	}
	baziChart.HourPillar = hourPillar

	jsonData, _ := json.Marshal(baziChart)
	slog.Info(string(jsonData))
	return &baziChart, nil
}
