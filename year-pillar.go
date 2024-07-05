package bazica

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

func GetYearPillar(path string, dateTime time.Time) (*YearPillar, error) {
	// convert solar date time to get lunar year
	lunarYear, err := GetLunarYear(path, dateTime)
	if err != nil {
		return nil, err
	}

	var yearPillar YearPillar
	yearPillar.Year = lunarYear
	stem := CalculateHeavenlyStem(lunarYear%10 - 3)
	yearPillar.HeavenlyStem = stem

	// calculate earthly branch (1900 is Rat year)
	branch := CalculateEarthlyBranch((lunarYear - 5) % 12)
	yearPillar.EarthlyBranch = branch

	return &yearPillar, nil
}

func GetLunarYear(path string, dateTime time.Time) (int, error) {
	// From 23:00 is new day (Rat hour)
	dateTime = dateTime.Add(time.Hour)

	solarYear := dateTime.Year()
	solarMonth := dateTime.Month()
	solarDay := dateTime.Day()

	lunarData := getNewYearData(path)

	var lunarYear int = solarYear
	// Check if the date is before or on the Lunar New Year of that solar year
	lunarNewYearDateStr, exists := lunarData.LunarNewYearDates[fmt.Sprintf("%d", solarYear)]
	if exists {
		lunarNewYearDateParts := strings.Split(lunarNewYearDateStr, "-")
		lunarNewYearMonth, _ := time.Parse("01", lunarNewYearDateParts[0])
		lunarNewYearDay, _ := time.Parse("02", lunarNewYearDateParts[1])

		if solarMonth < lunarNewYearMonth.Month() ||
			(solarMonth == lunarNewYearMonth.Month() && solarDay < lunarNewYearDay.Day()) {
			lunarYear--
		}
	}

	return lunarYear, nil
}

func getNewYearData(path string) *LunarNewYearData {
	// Assuming the combined JSON file is named "solar-term.json"
	prefix := PrefixPath
	if path != "" {
		prefix = path
	}

	fileToRead := prefix + "data/lunar-new-year.json"

	data, err := os.ReadFile(fileToRead)

	// Use slog directly from the beginning
	if err != nil {
		slog.Error("Error reading combined JSON", "error", err)
		return nil
	}

	var lunarData *LunarNewYearData
	err = json.Unmarshal(data, &lunarData)
	if err != nil {
		slog.Error("Error unmarshalling combined JSON", "error", err)
		return nil
	}

	// Check if combinedData is nil (optional)
	if lunarData == nil {
		slog.Warn("combinedData is nil. Check JSON structure and struct definition.")
		return nil
	} else {
		return lunarData
	}
}
