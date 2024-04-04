package bazica

import (
	svc "bazica/internal/datasvc"
	"bazica/models"
)

func CalTest() {

	svc.ReadJSONFile()
	tigerValue := models.Tiger
	println(tigerValue)

}
