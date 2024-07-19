package ultis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tommitoan/bazica/model"
	"log/slog"
	"os"
	"time"
)

func GetSolarTerm(path string, dateTime time.Time) (string, error) {
	// Extract the year
	yearStr := fmt.Sprint(dateTime.Year())
	result, err := GetSolarTermsByYear(yearStr, path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Use the correctly parsed time object when calling findSolarTerm
	term, passed, err := findSolarTerm(dateTime.Format("2006-01-02 15:04:05.999999999-07:00"), result)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error finding solar term: %v", err))
	} else {
		fmt.Println(passed)
		return term, nil
	}
}

func findSolarTerm(inputTime string, data model.SolarTermYear) (string, int, error) {
	t, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", inputTime)
	if err != nil {
		return "", 0, fmt.Errorf("invalid input time format: %v", err)
	}
	slog.Info("Input Time", "time", t)

	termList := []struct {
		name string
		time time.Time
	}{
		{"minor_cold", mustParseTime(data.MinorCold)},
		{"major_cold", mustParseTime(data.MajorCold)},
		{"start_of_spring", mustParseTime(data.StartOfSpring)},
		{"spring_showers", mustParseTime(data.SpringShowers)},
		{"awakening_of_insects", mustParseTime(data.AwakeningOfInsects)},
		{"spring_equinox", mustParseTime(data.SpringEquinox)},
		{"pure_brightness", mustParseTime(data.PureBrightness)},
		{"grain_rain", mustParseTime(data.GrainRain)},
		{"start_of_summer", mustParseTime(data.StartOfSummer)},
		{"grain_buds", mustParseTime(data.GrainBuds)},
		{"grain_in_ear", mustParseTime(data.GrainInEar)},
		{"summer_solstice", mustParseTime(data.SummerSolstice)},
		{"minor_heat", mustParseTime(data.MinorHeat)},
		{"major_heat", mustParseTime(data.MajorHeat)},
		{"start_of_autumn", mustParseTime(data.StartOfAutumn)},
		{"end_of_heat", mustParseTime(data.EndOfHeat)},
		{"white_dew", mustParseTime(data.WhiteDew)},
		{"autumn_equinox", mustParseTime(data.AutumnEquinox)},
		{"cold_dew", mustParseTime(data.ColdDew)},
		{"frost", mustParseTime(data.Frost)},
		{"start_of_winter", mustParseTime(data.StartOfWinter)},
		{"minor_snow", mustParseTime(data.MinorSnow)},
		{"major_snow", mustParseTime(data.MajorSnow)},
		{"winter_solstice", mustParseTime(data.WinterSolstice)},
	}

	var previousTermName string
	for i, term := range termList {
		if t.After(term.time) && (i+1 == len(termList) || t.Before(termList[i+1].time)) {
			timePassedMinutes := 0
			if i > 0 {
				timePassedMinutes = int(t.Sub(termList[i-1].time).Minutes())
			}
			slog.Info("Solar Term Found", "term", term.name, "timePassedMinutes", timePassedMinutes)
			return term.name, timePassedMinutes, nil // Return the name and time passed if the time falls within this term
		}
		if i > 0 {
			previousTermName = termList[i-1].name // Keep track of the previous term (except for the first term)
		}
	}
	// Handle the case where input time is before the first term in the list
	slog.Warn("Input time is before the first solar term", "previousTermName", previousTermName)
	return previousTermName, 0, nil // Return the last term of the previous year (or empty if there's none) and time passed 0
}

// Helper function to parse time with timezone offset
func mustParseTime(timeStr string) time.Time {
	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation("2006-01-02 15:04:05.999999999-07:00", timeStr, loc)
	if err != nil {
		panic("invalid time format in solar term data: " + err.Error())
	}
	return t
}

var PrefixPath string

func GetSolarTermsByYear(year string, path ...string) (model.SolarTermYear, error) {
	var prefix string
	if len(path) != 0 {
		prefix = path[0]
	}
	data := getSolarTermData(prefix)
	if data == nil {
		return model.SolarTermYear{}, errors.New("Cannot find solar term data for year " + year)
	}
	return data[year].Data, nil
}

func getSolarTermData(path string) map[string]model.CombinedData {
	// Assuming the combined JSON file is named "solar-term.json"
	prefix := PrefixPath
	if path != "" {
		prefix = path
	}

	fileToRead := prefix + "data/solar-term.json"

	data, err := os.ReadFile(fileToRead)

	// Use slog directly from the beginning
	if err != nil {
		slog.Error("Error reading combined JSON", "error", err)
		return nil
	}

	err = json.Unmarshal(data, &model.TempCombinedData)
	if err != nil {
		slog.Error("Error unmarshalling combined JSON", "error", err)
		return nil
	}

	// Check if combinedData is nil (optional)
	if model.TempCombinedData == nil {
		slog.Warn("combinedData is nil. Check JSON structure and struct definition.")
		return nil
	} else {
		return model.TempCombinedData
	}
}
