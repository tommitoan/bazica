package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type CombinedData struct {
	Year string            `json:"year"` // Year extracted from filename
	Data map[string]string `json:"data"` // Data for the specific year
}

func main() {
	// Define the output filename for combined data
	outputFile := "combined-6.json"

	// Define a slice to hold the combined data (for each year)
	combinedData := make(map[string]CombinedData)

	// Loop through the JSON filenames (1900.json to 2100.json)
	for i := 1900; i <= 2100; i++ {
		// Create the filename based on the loop counter
		filename := "svc/datasvc/solar-terms-data/" + strconv.Itoa(i) + ".json"

		// Read the content of the JSON file
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading file:", filename, err)
			continue
		}

		// Unmarshal the JSON data into a map
		var jsonData map[string]string
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", filename, err)
			continue
		}

		// Extract the year from the filename
		year := strconv.Itoa(i)

		// Add data for the year to the combined map
		combinedData[year] = CombinedData{
			Year: year,
			Data: jsonData,
		}
	}

	// Marshal the combined data back to JSON format
	jsonData, err := json.Marshal(combinedData)
	if err != nil {
		fmt.Println("Error marshalling combined data:", err)
		return
	}

	// Write the combined JSON data to the output file
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing combined JSON:", err)
		return
	}

	fmt.Println("Successfully wrote combined data to", outputFile)

	// Example: Get data for year 2000
	year := "2000"
	yearData, ok := combinedData[year]

	if ok {
		fmt.Println("Data for year", year, ":")
		fmt.Println("  Year:", yearData.Year)
		fmt.Println("  Data:", yearData.Data)
	} else {
		fmt.Println("Year", year, "not found in combined data.")
	}
}
