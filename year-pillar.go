package bazica

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

func GetYearPillar(path, dateStr string) (*YearPillar, error) {
	// convert solar date time to get lunar year
	lunarYear, err := GetLunarYear(path, dateStr)
	if err != nil {
		return nil, err
	}

	var yearPillar YearPillar
	yearPillar.Year = lunarYear
	stem := CalculateHeavenlyStem(lunarYear%10 - 3)
	yearPillar.HeavenlyStem = stem

	// calculate earthly branch (1900 is Rat year)
	switch (lunarYear - 5) % 12 {
	case RatValue:
		yearPillar.EarthlyBranch.Name = Rat
		yearPillar.EarthlyBranch.Value = RatValue
	case OxValue:
		yearPillar.EarthlyBranch.Name = Ox
		yearPillar.EarthlyBranch.Value = OxValue
	case TigerValue:
		yearPillar.EarthlyBranch.Name = Tiger
		yearPillar.EarthlyBranch.Value = TigerValue
	case RabbitValue:
		yearPillar.EarthlyBranch.Name = Rabbit
		yearPillar.EarthlyBranch.Value = RabbitValue
	case DragonValue:
		yearPillar.EarthlyBranch.Name = Dragon
		yearPillar.EarthlyBranch.Value = DragonValue
	case SnakeValue:
		yearPillar.EarthlyBranch.Name = Snake
		yearPillar.EarthlyBranch.Value = SnakeValue
	case HorseValue:
		yearPillar.EarthlyBranch.Name = Horse
		yearPillar.EarthlyBranch.Value = HorseValue
	case GoatValue:
		yearPillar.EarthlyBranch.Name = Goat
		yearPillar.EarthlyBranch.Value = GoatValue
	case MonkeyValue:
		yearPillar.EarthlyBranch.Name = Monkey
		yearPillar.EarthlyBranch.Value = MonkeyValue
	case RoosterValue:
		yearPillar.EarthlyBranch.Name = Rooster
		yearPillar.EarthlyBranch.Value = RoosterValue
	case DogValue:
		yearPillar.EarthlyBranch.Name = Dog
		yearPillar.EarthlyBranch.Value = DogValue
	case PigValue:
		yearPillar.EarthlyBranch.Name = Pig
		yearPillar.EarthlyBranch.Value = PigValue
	}

	return &yearPillar, nil
}

func GetLunarYear(path, dateStr string) (int, error) {
	solarDate, err := GetDateFormat(dateStr)
	if err != nil {
		return 0, fmt.Errorf("error parsing date string: %v", err)
	}

	// Handle 00:00 to 00:59 still on Previous day (Rat hour)
	solarDate = solarDate.Add(-time.Hour)

	solarYear := solarDate.Year()
	solarMonth := solarDate.Month()
	solarDay := solarDate.Day()

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
