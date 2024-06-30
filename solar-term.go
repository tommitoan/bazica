package bazica

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
)

var PrefixPath string

func GetSolarTermsByYear(year string, path ...string) (SolarTermYear, error) {
	var prefix string
	if len(path) != 0 {
		prefix = path[0]
	}
	data := getSolarTermData(prefix)
	if data == nil {
		return SolarTermYear{}, errors.New("Cannot find solar term data for year " + year)
	}
	return data[year].Data, nil
}

func getSolarTermData(path string) map[string]CombinedData {
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
