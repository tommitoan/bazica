package bazica

import (
	"github.com/tommitoan/bazica/model"
	"time"
)

func GetBaziChart(path string, dateTime time.Time, loc *time.Location) (*model.BaziChart, error) {
	var baziChart model.BaziChart

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

	return &baziChart, nil
}