package analysis

import (
	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

// patternCandidates records which month-branch hidden stems are exposed among
// the three non-day stems. It does not resolve combinations or special patterns.
func patternCandidates(pillars [4]natalPillar) []model.PatternCandidate {
	hidden := hiddenStemsOf(pillars[1].branch.Value)
	dayStem := pillars[2].stem.Value
	result := make([]model.PatternCandidate, 0, len(hidden))
	for rank, stem := range hidden {
		exposed := make([]string, 0, 3)
		for _, i := range []int{0, 1, 3} {
			if pillars[i].stem.Value == stem {
				exposed = append(exposed, pillarKeys[i])
			}
		}
		if len(exposed) == 0 {
			continue
		}
		god := tenGod(dayStem, stem)
		result = append(result, model.PatternCandidate{
			TenGod: god, HiddenStem: utils.CalculateHeavenlyStem(stem),
			QiRank: rank + 1, ExposedAt: exposed, Basis: "month_hidden_stem_exposed",
			Regular: god.Code != "ten_god.friend" && god.Code != "ten_god.rob_wealth",
		})
	}
	if len(result) > 0 {
		return result
	}
	stem := hidden[0]
	god := tenGod(dayStem, stem)
	return []model.PatternCandidate{{
		TenGod: god, HiddenStem: utils.CalculateHeavenlyStem(stem),
		QiRank: 1, ExposedAt: []string{}, Basis: "month_main_qi_unexposed",
		Regular: god.Code != "ten_god.friend" && god.Code != "ten_god.rob_wealth",
	}}
}
