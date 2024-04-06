package datasvc

import (
	"bazica/internal/toerr"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

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

func GetSolarTermsByYear(year string) (*SolarTerms, error) {

	fileToRead := "internal/datasvc/solar-terms-data/" + year + ".json"
	// Open file
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
	// Open the JSON file
	//decoder := json.NewDecoder(file)
	//fmt.Println(decoder)
	//fmt.Println("hello")
	//
	//// Create an instance of the SolarTerms struct
	//var solarTerms SolarTerms
	//
	//// Unmarshal the JSON data into the struct
	//err = json.Unmarshal(data, &solarTerms)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Access the data from the struct
	//log.Println("Minor Cold:", solarTerms.MinorCold)
	//log.Println("Major Cold:", solarTerms.MajorCold)
}
