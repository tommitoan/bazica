package bazica

import (
	svc "bazica/internal/datasvc"
	"bazica/models"
)

const (
	YearToRead = "1950"
	FileToRead = "internal/datasvc/solar-terms-data/" + YearToRead + ".json"
)

func CalTest() {

	svc.ReadJSONFile(FileToRead)
	tigerValue := models.Tiger
	println(tigerValue)

}
