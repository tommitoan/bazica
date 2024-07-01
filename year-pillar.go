package bazica

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

func GetYearPillar(path, dateStr string) (string, error) {
	dateParts := strings.Split(dateStr, "T") // Split date and time (if present)
	if len(dateParts) == 0 {
		return "", fmt.Errorf("invalid date string format")
	}
	dateStr = dateParts[0] // Extract only the date part (YYYY-MM-DD)

	solarDate, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", dateStr)
	if err != nil {
		return "", fmt.Errorf("error parsing date string: %v", err)
	}
	solarDate = solarDate.Add(-time.Hour) // Handle 00:00 to 00:59 still on Previous day (Rat hour)

	solarYear := solarDate.Year()
	solarMonth := solarDate.Month()
	solarDay := solarDate.Day()

	lunarData := getNewYearData(path)

	// Check if the date is before or on the Lunar New Year of that solar year
	lunarNewYearDateStr, exists := lunarData.LunarNewYearDates[fmt.Sprintf("%d", solarYear)]
	if exists {
		lunarNewYearDateParts := strings.Split(lunarNewYearDateStr, "-")
		lunarNewYearMonth, _ := time.Parse("01", lunarNewYearDateParts[0])
		lunarNewYearDay, _ := time.Parse("02", lunarNewYearDateParts[1])

		if solarMonth < lunarNewYearMonth.Month() ||
			(solarMonth == lunarNewYearMonth.Month() && solarDay <= lunarNewYearDay.Day()) {
			return fmt.Sprintf("%d", solarYear-1), nil
		}
	}

	fmt.Println(lunarNewYearDateStr)
	return "", nil
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

	slog.Info("Successfully unmarshalled data!")

	// Check if combinedData is nil (optional)
	if lunarData == nil {
		slog.Warn("combinedData is nil. Check JSON structure and struct definition.")
		return nil
	} else {
		slog.Info("combinedData contains data!")
		return lunarData
	}
}
