package main

import (
	solar_terms_data "bazica/internal/datasvc/solar-terms-data"
	"fmt"
)

func main() {
	path := "internal/datasvc/solar-terms-data/"
	yearToCheck := "1995"

	slt, _ := solar_terms_data.GetSolarTermsByYear(path, yearToCheck)

	fmt.Println(slt.AutumnEquinox)
	//println(models.Tiger)
	//println(models.Tiger.String())
	//fmt.Println(solartermsvc.GetSolarTermByYear("1900"))
}
