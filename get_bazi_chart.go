package bazica

import (
	"github.com/tommitoan/bazica/internal/fourpillars"
	"github.com/tommitoan/bazica/internal/luckpillars"
	"github.com/tommitoan/bazica/internal/ultis"
	"github.com/tommitoan/bazica/model"
	"time"
)

func GetBaziChart(dateTime time.Time, loc *time.Location, gender int) (*model.BaziChart, error) {
	var baziChart model.BaziChart

	fourPillar, passed, remaining, err := fourpillars.GetFourPillars(dateTime, loc)
	if err != nil {
		return nil, err
	}
	baziChart.FourPillar = ultis.GetLifeCycleFromFourPillar(fourPillar)

	lucksPillar, err := luckpillars.GetLuckPillars(fourPillar, gender, passed, remaining, dateTime)
	if err != nil {
		return nil, err
	}
	baziChart.LuckPillars = lucksPillar

	return &baziChart, nil
}
