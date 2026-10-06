package fourpillars

import (
	"fmt"
	"time"

	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

// GetYearReference returns the display-only year data of a birth: the Lichun of
// its civil year in the birth zone and the year of the lunar calendar. Nothing
// in the chart is calculated from it. The lunar-calendar year follows the civil
// date, because the lunar calendar changes year at midnight, unlike the Ba-zi day.
func GetYearReference(dateTime time.Time, yearPillar *model.YearPillar) (*model.YearReference, error) {
	lichun, err := utils.StartOfSpring(dateTime.Year())
	if err != nil {
		return nil, err
	}
	lunarYear, err := GetLunarYear(dateTime)
	if err != nil {
		return nil, err
	}
	stem, branch := yearStemBranch(lunarYear)

	reference := &model.YearReference{
		Lichun: model.LichunReference{
			Time:              lichun.In(dateTime.Location()),
			BirthBeforeLichun: dateTime.Before(lichun),
		},
		LunarYear: model.LunarYearReference{
			Year:          lunarYear,
			HeavenlyStem:  stem,
			EarthlyBranch: branch,
			GanZhi:        utils.GanZhiOf(stem, branch),
		},
		Differs: stem.Value != yearPillar.HeavenlyStem.Value || branch.Value != yearPillar.EarthlyBranch.Value,
	}
	lunarData, err := loadLunarNewYear()
	if err != nil {
		return nil, err
	}
	if monthDay, ok := lunarData.LunarNewYearDates[fmt.Sprintf("%d", lunarYear)]; ok {
		reference.LunarYear.NewYearDate = fmt.Sprintf("%04d-%s", lunarYear, monthDay)
	}
	return reference, nil
}
