package fourpillars

import (
	"encoding/json"
	"fmt"
	"github.com/tommitoan/bazica/data"
	"github.com/tommitoan/bazica/internal/ultis"
	"github.com/tommitoan/bazica/model"
	"strings"
	"sync"
	"time"
)

func GetYearPillar(dateTime time.Time) (*model.YearPillar, error) {
	// convert solar date time to get lunar year
	lunarYear, err := GetLunarYear(dateTime)
	if err != nil {
		return nil, err
	}

	var yearPillar model.YearPillar
	yearPillar.Year = dateTime.Year()

	stemValue := lunarYear%10 - 3
	if stemValue < 1 {
		stemValue = stemValue + 10
	}
	stem := ultis.CalculateHeavenlyStem(stemValue)
	yearPillar.HeavenlyStem = stem

	// calculate earthly branch (1900 is Rat year)
	branchValue := (lunarYear - 5) % 12
	if branchValue < 1 {
		branchValue = branchValue + 12
	}
	branch := ultis.CalculateEarthlyBranch(branchValue)
	yearPillar.EarthlyBranch = branch

	return &yearPillar, nil
}

func GetLunarYear(dateTime time.Time) (int, error) {
	// From 23:00 is new day (Rat hour)
	dateTime = dateTime.Add(time.Hour)

	solarYear := dateTime.Year()
	solarMonth := dateTime.Month()
	solarDay := dateTime.Day()

	lunarData, err := loadLunarNewYear()
	if err != nil {
		return 0, err
	}

	var lunarYear int = solarYear
	// Check if the date is before or on the Lunar New Year of that solar year
	lunarNewYearDateStr, exists := lunarData.LunarNewYearDates[fmt.Sprintf("%d", solarYear)]
	if !exists {
		// Jan 1 of the first year without data still precedes that year's
		// Lunar New Year, so it belongs to the last covered lunar year.
		if _, prevExists := lunarData.LunarNewYearDates[fmt.Sprintf("%d", solarYear-1)]; prevExists && solarMonth == time.January && solarDay == 1 {
			return solarYear - 1, nil
		}
		return 0, fmt.Errorf("%w: no Lunar New Year data for %d", model.ErrDateOutOfRange, solarYear)
	}
	lunarNewYearDateParts := strings.Split(lunarNewYearDateStr, "-")
	lunarNewYearMonth, _ := time.Parse("01", lunarNewYearDateParts[0])
	lunarNewYearDay, _ := time.Parse("02", lunarNewYearDateParts[1])

	if solarMonth < lunarNewYearMonth.Month() ||
		(solarMonth == lunarNewYearMonth.Month() && solarDay < lunarNewYearDay.Day()) {
		lunarYear--
	}

	return lunarYear, nil
}

// loadLunarNewYear decodes the embedded table once; the result is read-only.
var loadLunarNewYear = sync.OnceValues(func() (*model.LunarNewYearData, error) {
	var lunarData model.LunarNewYearData
	if err := json.Unmarshal(data.LunarNewYear, &lunarData); err != nil {
		return nil, fmt.Errorf("%w: lunar-new-year.json: %v", model.ErrDataUnavailable, err)
	}
	if len(lunarData.LunarNewYearDates) == 0 {
		return nil, fmt.Errorf("%w: lunar-new-year.json is empty", model.ErrDataUnavailable)
	}
	return &lunarData, nil
})
