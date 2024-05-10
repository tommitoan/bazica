package main

//
//import (
//	"encoding/json"
//	"fmt"
//	"os"
//	"strconv"
//)
//
//type CombinedData struct {
//	Year string            `json:"year"` // Year extracted from filename
//	Data map[string]string `json:"data"` // Data for the specific year
//}
//
//func main() {
//	// Define the output filename for combined data
//	outputFile := "combined-5.json"
//
//	// Define a slice to hold the combined data (for each year)
//	var combinedData []CombinedData
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
//		// Create a CombinedData struct for this year
//		yearData := CombinedData{
//			Year: strconv.Itoa(i),
//			Data: jsonData,
//		}
//
//		// Append the data to the combinedData slice
//		combinedData = append(combinedData, yearData)
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
//	fmt.Println("Successfully wrote combined data to", outputFile)
//}
