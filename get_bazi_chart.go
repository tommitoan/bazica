package bazica

import (
	"time"

	"github.com/tommitoan/bazica/internal/fourpillars"
	"github.com/tommitoan/bazica/internal/luckpillars"
	"github.com/tommitoan/bazica/internal/ultis"
	"github.com/tommitoan/bazica/model"
)

func GetBaziChart(dateTime time.Time, loc *time.Location, gender int, prefixPath ...string) (*model.BaziChart, error) {
	var baziChart model.BaziChart

	var path string
	if len(prefixPath) != 0 {
		path = prefixPath[0]
	}

	fourPillar, passed, remaining, err := fourpillars.GetFourPillars(dateTime, loc, path)
	if err != nil {
		return nil, err
	}
	baziChart.FourPillar = ultis.GetLifeCycleFromFourPillar(fourPillar)
	baziChart.ElementScore = ultis.GetElementScoreFromFourPillar(baziChart.FourPillar)

	lucksPillar, err := luckpillars.GetLuckPillars(fourPillar, gender, passed, remaining, dateTime, path)
	if err != nil {
		return nil, err
	}
	baziChart.LuckPillars = lucksPillar

	return &baziChart, nil
}
