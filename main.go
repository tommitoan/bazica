package main

import (
	"bazica/internal/datasvc"
	"fmt"
)

func main() {
	slt, _ := datasvc.GetSolarTermsByYear("2101")

	fmt.Println(slt.AutumnEquinox)
	//println(models.Tiger)
	//println(models.Tiger.String())
	//fmt.Println(solartermsvc.GetSolarTermByYear("1900"))
}
