package fourpillars

import (
	"encoding/json"
	"fmt"
	"github.com/tommitoan/bazica/v2/data"
	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
	"strings"
	"sync"
	"time"
)

// GetYearPillar returns the Ba-zi year pillar, which changes at the instant of
// Lichun (Start of Spring), the same solar-term clock that drives the month
// pillar. YearPillar.Year is the civil year of the birth.
func GetYearPillar(dateTime time.Time) (*model.YearPillar, error) {
	baziYear, err := GetBaziYear(dateTime)
	if err != nil {
		return nil, err
	}

	var yearPillar model.YearPillar
	yearPillar.Year = dateTime.Year()
	yearPillar.HeavenlyStem, yearPillar.EarthlyBranch = yearStemBranch(baziYear)
	return &yearPillar, nil
}

// yearStemBranch returns the stem and branch of a Gregorian year (1984 was a
// Jia-Zi year, 1900 a Geng-Zi year).
func yearStemBranch(year int) (model.HeavenlyStem, model.EarthlyBranch) {
	stemValue := year%10 - 3
	if stemValue < 1 {
		stemValue = stemValue + 10
	}
	branchValue := (year - 5) % 12
	if branchValue < 1 {
		branchValue = branchValue + 12
	}
	return utils.CalculateHeavenlyStem(stemValue), utils.CalculateEarthlyBranch(branchValue)
}

// GetBaziYear returns the Gregorian year the year pillar stands for: the year
// of the latest Lichun at dateTime. The comparison uses the instant, like the
// month pillar, so the late-Rat-hour day rule does not apply to it.
func GetBaziYear(dateTime time.Time) (int, error) {
	lichun, err := utils.StartOfSpring(dateTime.Year())
	if err != nil {
		return 0, err
	}
	if dateTime.Before(lichun) {
		return dateTime.Year() - 1, nil
	}
	return dateTime.Year(), nil
}

// GetLunarYear returns the Gregorian year in which the lunar-calendar year of dateTime
// begins. It is reference data only: the year pillar follows Lichun (GetBaziYear).
// The lunar calendar changes year at midnight, so the civil date in the birth zone is
// compared with the Lunar New Year date; the 23:00 rule of the day pillar does not apply.
func GetLunarYear(dateTime time.Time) (int, error) {
	solarYear, solarMonth, solarDay := dateTime.Date()

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

// SupportedYears returns the first and last calendar year of a birth date the
// embedded tables cover. A year y is covered when the solar-term table holds
// y-1, y and y+1 and the Lunar New Year table holds y, so the range follows the
// data and never has to be repeated as a constant.
func SupportedYears() (first, last int, err error) {
	solarFirst, solarLast, err := utils.SolarTermYears()
	if err != nil {
		return 0, 0, err
	}
	lunarData, err := loadLunarNewYear()
	if err != nil {
		return 0, 0, err
	}
	lunarFirst, lunarLast := 0, 0
	for key := range lunarData.LunarNewYearDates {
		var year int
		if _, scanErr := fmt.Sscanf(key, "%d", &year); scanErr != nil {
			return 0, 0, fmt.Errorf("%w: lunar-new-year.json: year %q", model.ErrDataUnavailable, key)
		}
		if lunarFirst == 0 || year < lunarFirst {
			lunarFirst = year
		}
		if year > lunarLast {
			lunarLast = year
		}
	}
	return max(solarFirst+1, lunarFirst), min(solarLast-1, lunarLast), nil
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
