package main

import (
	"fmt"
	"github.com/tommitoan/bazica/internal/datasvc"
)

func main() {
	path := "internal/datasvc/solar-terms-data/"
	yearToCheck := "1995"

	slt, _ := datasvc.GetSolarTermsByYear(path, yearToCheck)

	fmt.Println(slt.AutumnEquinox)
	//println(models.Tiger)
	//println(models.Tiger.String())
	//fmt.Println(solartermsvc.GetSolarTermByYear("1900"))
}
