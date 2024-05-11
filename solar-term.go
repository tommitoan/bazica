package bazica

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
)

const PrefixPathSolarTerm = ""

type SolarTermYear struct {
	MinorCold          string `json:"minor_cold"`
	MajorCold          string `json:"major_cold"`
	StartOfSpring      string `json:"start_of_spring"`
	SpringShowers      string `json:"spring_showers"`
	AwakeningOfInsects string `json:"awakening_of_insects"`
	SpringEquinox      string `json:"spring_equinox"`
	PureBrightness     string `json:"pure_brightness"`
	GrainRain          string `json:"grain_rain"`
	StartOfSummer      string `json:"start_of_summer"`
	GrainBuds          string `json:"grain_buds"`
	GrainInEar         string `json:"grain_in_ear"`
	SummerSolstice     string `json:"summer_solstice"`
	MinorHeat          string `json:"minor_heat"`
	MajorHeat          string `json:"major_heat"`
	StartOfAutumn      string `json:"start_of_autumn"`
	EndOfHeat          string `json:"end_of_heat"`
	WhiteDew           string `json:"white_dew"`
	AutumnEquinox      string `json:"autumn_equinox"`
	ColdDew            string `json:"cold_dew"`
	Frost              string `json:"frost"`
	StartOfWinter      string `json:"start_of_winter"`
	MinorSnow          string `json:"minor_snow"`
	MajorSnow          string `json:"major_snow"`
	WinterSolstice     string `json:"winter_solstice"`
}

type CombinedData struct {
	Year string        `json:"year"` // Year extracted from filename
	Data SolarTermYear `json:"data"` // Data for the specific year
}

var combinedData map[string]CombinedData

func GetSolarTermsByYear(year, path string) (SolarTermYear, error) {
	data := getSolarTermData(path)
	if data == nil {
		return SolarTermYear{}, errors.New("Cannot find solar term data for year " + year)
	}
	return data[year].Data, nil
}

func getSolarTermData(path string) map[string]CombinedData {
	// Assuming the combined JSON file is named "solar-term.json"
	prefix := PrefixPathSolarTerm
	if path != "" {
		prefix = path
	}

	fileToRead := prefix + "internal/solar-term.json"
	data, err := os.ReadFile(fileToRead)

	// Use slog directly from the beginning
	if err != nil {
		slog.Error("Error reading combined JSON", "error", err)
		return nil
	}

	err = json.Unmarshal(data, &combinedData)
	if err != nil {
		slog.Error("Error unmarshalling combined JSON", "error", err)
		return nil
	}

	slog.Info("Successfully unmarshalled data!")

	// Check if combinedData is nil (optional)
	if combinedData == nil {
		slog.Warn("combinedData is nil. Check JSON structure and struct definition.")
		return nil
	} else {
		slog.Info("combinedData contains data!")
		return combinedData
	}
}
