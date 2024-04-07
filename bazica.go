package bazica

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/tommitoan/bazica/svc/datasvc"
	toerr "github.com/tommitoan/bazica/svc/toerr"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	PrefixLocal  = "svc/datasvc/solar-terms-data/"
	PrefixVendor = "vendor/github.com/tommitoan/bazica/svc/datasvc/solar-terms-data/"
)

//go:embed svc/datasvc/solar-terms-data/*
var dataDir embed.FS

type BazicaSvc struct {
	DataSvc *datasvc.DataSvc
	Solar   *datasvc.SolarTerms
	dataDir embed.FS
}

func (bzc *BazicaSvc) GetSolarTermsByYear(prefix, year string) (*datasvc.SolarTerms, error) {
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
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "Error when open file")
		log.Fatal(err)
	} else if file == nil {
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "Wrong path")
	}

	defer file.Close()

	// Read as byte array
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "Can't read file")
	}
	var solarTerm datasvc.SolarTerms
	json.Unmarshal(byteValue, &solarTerm)

	return &solarTerm, nil
}

func (bzc *BazicaSvc) GetSolarTermsByYearV3(year string) (*datasvc.SolarTerms, error) {
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
	fileToRead := year + ".json"

	fmt.Println(dataDir)
	files, err := dataDir.ReadDir("svc/datasvc/solar-terms-data/") // Replace with desired directory
	if err != nil {
		return nil, err
	}

	fmt.Println("Available files:", files)

	file, err := dataDir.Open(fileToRead)
	if err != nil {
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "Error when open file")
		log.Fatal(err)
	} else if file == nil {
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "Wrong path")
	}

	defer file.Close()

	// Read as byte array
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, toerr.NewValidationError(http.StatusInternalServerError, "Can't read file")
	}
	var solarTerm datasvc.SolarTerms
	json.Unmarshal(byteValue, &solarTerm)

	return &solarTerm, nil
}
