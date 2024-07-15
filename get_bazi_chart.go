package bazica

import (
	"github.com/tommitoan/bazica/internal/fourpillars"
	"github.com/tommitoan/bazica/internal/luckpillars"
	"github.com/tommitoan/bazica/model"
	"time"
)

func GetBaziChart(dateTime time.Time, loc *time.Location, prefixPath ...string) (*model.BaziChart, error) {
	// TODO: ask gender for calculating luck pillars
	var gender int = 1

	var baziChart model.BaziChart

	var path string
	if len(prefixPath) != 0 {
		path = prefixPath[0]
	}

	fourPillar, err := fourpillars.GetFourPillars(dateTime, loc, path)
	if err != nil {
		return nil, err
	}
	baziChart.FourPillar = fourPillar

	lucksPillar, err := luckpillars.GetLuckPillars(fourPillar, gender, path)
	if err != nil {
		return nil, err
	}
	baziChart.LuckPillars = lucksPillar

	return &baziChart, nil
}
