package bazica

import (
	svc "bazica/internal/datasvc/solar-terms-data"
	"fmt"
)

const (
	YearToRead = "1950"
	FileToRead = "internal/datasvc/solar-terms-data/" + YearToRead + ".json"
)

func ReadSolarTerms() {

	abc, _ := svc.GetSolarTermsByYear("1950")

	fmt.Println(abc.ColdDew, abc.AutumnEquinox)
	tigerValue := models.Tiger
	println(tigerValue)

}
