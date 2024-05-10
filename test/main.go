package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("combined.json")

	if err != nil {
		fmt.Println("Error reading combined JSON:", err)
		return
	}
	fmt.Println(string(data))
}
