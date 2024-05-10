package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type CombinedData struct {
	Year string            `json:"year"` // Year extracted from filename
	Data map[string]string `json:"data"` // Data for the specific year
}

var combinedData map[string]CombinedData

func main() {
	// Assuming the combined JSON file is named "combined.json"

	data, err := os.ReadFile("test/solar-term-v3.json")
	if err != nil {
		fmt.Println("Error reading combined JSON:", err)
		return
	}

	err = json.Unmarshal(data, &combinedData)
	if err != nil {
		fmt.Println("Error unmarshalling combined JSON:", err)
		return
	}

	fmt.Println("Successfully unmarshalled data!")

	// Check if combinedData is nil (optional)
	if combinedData == nil {
		fmt.Println("combinedData is nil. Check JSON structure and struct definition.")
	} else {
		fmt.Println("combinedData contains data!")
	}
}
