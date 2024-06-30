package bazica

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

func DetectSolarTerm(path, input string) (string, error) {
	// Parse the timestamp string (updated layout)
	dt, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", input)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return "", err
	}

	// Extract the year
	yearStr := fmt.Sprint(dt.Year())
	result, err := GetSolarTermsByYear(yearStr, path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Use the correctly parsed time object when calling findSolarTerm
	term, err := findSolarTerm(dt.Format("2006-01-02 15:04:05.999999999-07:00"), result)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error finding solar term: %v", err))
	} else {
		return term, nil
	}
}

func findSolarTerm(inputTime string, data SolarTermYear) (string, error) {
	t, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", inputTime)
	if err != nil {
		return "", fmt.Errorf("invalid input time format: %v", err)
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
		slog.Info("Checking Solar Term", "term", term.name, "termTime", term.time)

		if t.After(term.time) && (i+1 == len(termList) || t.Before(termList[i+1].time)) {
			slog.Info("Solar Term Found", "term", term.name)
			return term.name, nil // Return the name if the time falls within this term
		}
		if i > 0 {
			previousTermName = termList[i-1].name // Keep track of the previous term (except for the first term)
			slog.Info("Updating Previous Term", "previousTermName", previousTermName)
		}
	}
	// Handle the case where input time is before the first term in the list
	slog.Warn("Input time is before the first solar term", "previousTermName", previousTermName) // Changed to Warning as this isn't an error per se
	return previousTermName, nil                                                                 // Return the last term of the previous year (or empty if there's none)
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
