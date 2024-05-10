package datasvc

import (
	"encoding/json"
	"fmt"
	toerr "github.com/tommitoan/bazica/svc/toerr"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type CombinedData struct {
	Year string     `json:"year"` // Year extracted from filename
	Data SolarTerms `json:"data"` // Data for the specific year
}

type DataSvc struct {
	Year  string
	Solar SolarTerms
}

type SolarTerms struct {
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

func (dsvc *DataSvc) GetSolarTermsByYearV2(prefix, year string) (*SolarTerms, error) {
	// Handle year from 1900 -> 2100 only
	i, err := strconv.Atoi(year)
	if err != nil {
		panic(err)
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "Something wrong")
	}
	if i < 1900 || i > 2100 {
		return nil, toerr.NewValidationError(http.StatusBadRequest, "Year not found")
	}

	// Open file
	fileToRead := prefix + year + ".json"
	file, err := os.Open(fileToRead)
	if err != nil {
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "File not found")
		log.Fatal(err)
	}
	defer file.Close()

	// Read as byte array
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "Can't read file")
	}
	var solarTerm SolarTerms
	json.Unmarshal(byteValue, &solarTerm)

	return &solarTerm, nil
}

func (dsvc *DataSvc) GetSolarTermsByYearV3(prefix, year string) (*SolarTerms, error) {
	// Assuming the combined JSON file is named "combined.json"
	data, err := os.ReadFile("combined.json")
	if err != nil {
		fmt.Println("Error reading combined JSON:", err)
		return nil, err
	}

	var combinedData map[string]CombinedData // Declare as map

	err = json.Unmarshal(data, &combinedData)
	if err != nil {
		fmt.Println("Error unmarshalling combined JSON:", err)
		return nil, err
	}

	fmt.Println("Successfully unmarshalled data!")

	// Get data for a specific year (example: 2000)
	yearData, ok := combinedData[year] // Check if key exists

	if ok {
		fmt.Println("Data for year", year, ":")
		fmt.Println("  Year:", yearData.Year)
		fmt.Println("  Data:", yearData.Data)
	} else {
		fmt.Println("Year", year, "not found in combined data.")
	}

	fmt.Println("Solar Terms for year:", yearData.Year)
	return nil, nil
}

var combinedData CombinedData
