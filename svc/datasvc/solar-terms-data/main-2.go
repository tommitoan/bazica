package main

//
//import (
//	"encoding/json"
//	"fmt"
//	"os"
//	"strconv"
//)
//
//func main() {
//	// Define the output filename
//	outputFile := "combined.json"
//
//	// Define a map to hold the combined data (year as string key, data as map)
//	combinedData := make(map[string]map[string]string)
//
//	// Loop through the JSON filenames (1900.json to 2100.json)
//	for i := 1900; i <= 2100; i++ {
//		// Create the filename based on the loop counter
//		filename := "svc/datasvc/solar-terms-data/" + strconv.Itoa(i) + ".json"
//
//		// Read the content of the JSON file
//		data, err := os.ReadFile(filename)
//		if err != nil {
//			fmt.Println("Error reading file:", filename, err)
//			continue
//		}
//
//		// Unmarshal the JSON data into a map
//		var jsonData map[string]string
//		err = json.Unmarshal(data, &jsonData)
//		if err != nil {
//			fmt.Println("Error unmarshalling JSON:", filename, err)
//			continue
//		}
//
//		// Extract the year from the filename
//		year := strconv.Itoa(i)
//
//		// Add data for the year to the combined map
//		combinedData[year] = jsonData
//	}
//
//	// Marshal the combined data back to JSON format
//	jsonData, err := json.Marshal(combinedData)
//	if err != nil {
//		fmt.Println("Error marshalling combined data:", err)
//		return
//	}
//
//	// Write the combined JSON data to the output file
//	err = os.WriteFile(outputFile, jsonData, 0644)
//	if err != nil {
//		fmt.Println("Error writing combined JSON:", err)
//		return
//	}
//
//	fmt.Println("Successfully combined JSON files into", outputFile)
//
//	// Example: Access data for a specific year (assuming data exists)
//	yearToAccess := "1900"
//	dataForYear, ok := combinedData[yearToAccess]
//	if !ok {
//		fmt.Println("Data not found for year:", yearToAccess)
//		return
//	}
//	fmt.Println("Data for year", yearToAccess, ":", dataForYear)
//}
