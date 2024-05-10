package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Define the output filename
	outputFile := "solar-term.json"

	// Define a slice to hold the combined data
	var combinedData []interface{}

	// Loop through the JSON files (1900.json to 2100.json)
	for i := 1900; i <= 2100; i++ {
		// Create the filename based on the loop counter
		filename := "svc/datasvc/solar-terms-data/" + strconv.Itoa(i) + ".json"

		// Read the content of the JSON file
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading file:", filename, err)
			continue
		}

		// Unmarshal the JSON data
		var jsonData map[string]string
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", filename, err)
			continue
		}

		// Extract the year from the filename
		year, err := strconv.Atoi(filename[29:33])
		if err != nil {
			fmt.Println("Error parsing year from filename:", filename, err)
			continue
		}

		// Add "year" key with extracted year value
		jsonData["year"] = strconv.Itoa(year)

		// Append the data to the combined slice
		combinedData = append(combinedData, jsonData)
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

	fmt.Println("Successfully combined JSON files into", outputFile)
}
