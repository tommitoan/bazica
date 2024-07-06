package bazica

import (
	"fmt"
	"github.com/tommitoan/bazica/internal/fourpillars"
	"github.com/tommitoan/bazica/model"
	"time"
)

func GetBaziChart(dateTime time.Time, loc *time.Location, prefixPath ...string) (*model.BaziChart, error) {
	var path string
	if len(prefixPath) != 0 {
		path = prefixPath[0]
	}

	fourPillar, err := fourpillars.GetFourPillars(dateTime, loc, path)
	if err != nil {
		return nil, err
	}
	fmt.Println(fourPillar)

	return nil, nil
}
