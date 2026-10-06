package analysis

import (
	"fmt"

	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

// attachLuck fills the Analysis field of every luck pillar of chart. Ages are
// nominal: the year of the year pillar counts as age 1, so a luck pillar that
// starts in the year after the year pillar's year begins at age 2.
func attachLuck(chart *model.BaziChart, pillars [4]natalPillar) error {
	if chart.LuckPillars == nil {
		return nil
	}
	dayStem := pillars[2].stem.Value
	natal := natalCells(pillars)
	baseYear := yearPillarYear(chart.FourPillar.YearPillar)

	for _, lp := range chart.LuckPillars.LuckPillars {
		if lp == nil {
			continue
		}
		stem, branch := lp.HeavenlyStem.Value, lp.EarthlyBranch.Value
		if !validStem(stem) || !validBranch(branch) {
			return fmt.Errorf("%w: luck pillar %d has stem %d, branch %d", errInvalidPillar, lp.Number, stem, branch)
		}
		nayin, ok := nayinTerms[lp.GanZhi.Name]
		if !ok {
			return fmt.Errorf("%w: no Nayin for %q", errInvalidPillar, lp.GanZhi.Name)
		}
		lp.Analysis = &model.LuckPillarAnalysis{
			AgeStart:             lp.YearStart - baseYear + 1,
			AgeEnd:               lp.YearEnd - baseYear + 1,
			TenGod:               tenGod(dayStem, stem),
			Nayin:                nayin,
			StemStageAtOwnBranch: stage(stem, branch),
			HiddenStems:          luckHiddenStemList(dayStem, branch),
			HeavenEarthClash:     heavenEarthClash(cell{stem, branch}, natal),
		}
	}
	return nil
}

// luckHiddenStemList lists the stems hidden in branch with their stage at the
// branch itself; luck pillars have no month branch of their own.
func luckHiddenStemList(dayStem, branch int) []model.LuckHiddenStem {
	values := hiddenStemsOf(branch)
	list := make([]model.LuckHiddenStem, 0, len(values))
	for _, h := range values {
		list = append(list, model.LuckHiddenStem{
			HeavenlyStem:     utils.CalculateHeavenlyStem(h),
			TenGod:           tenGod(dayStem, h),
			StageAtOwnBranch: stage(h, branch),
		})
	}
	return list
}
