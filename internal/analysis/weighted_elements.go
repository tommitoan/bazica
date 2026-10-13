package analysis

import (
	"sort"

	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

var elementKeys = [5]string{"wood", "fire", "earth", "metal", "water"}
var pillarKeys = [4]string{"year", "month", "day", "hour"}

// The weights follow the ordered hidden-stem table: main qi first.
var hiddenWeights = [4][3]int{{}, {10}, {7, 3}, {6, 3, 1}}

func weightedLevel(tenths int) int {
	switch {
	case tenths < 5:
		return 0
	case tenths < 15:
		return 1
	case tenths < 25:
		return 2
	case tenths < 35:
		return 3
	default:
		return 4
	}
}

func weightedElements(pillars [4]natalPillar) model.WeightedElements {
	type bucket struct {
		score, dayMaster int
		gods             map[string]int
		contributions    []model.ElementContribution
	}
	var buckets [5]bucket
	for i := range buckets {
		buckets[i].gods = make(map[string]int)
		buckets[i].contributions = make([]model.ElementContribution, 0)
	}
	dayStem := pillars[2].stem.Value
	add := func(pillar, source string, stem int, branch *model.EarthlyBranch, tenths int, dayMaster bool) {
		element := stemElement(stem)
		b := &buckets[element]
		var god *model.LocalizedTerm
		if dayMaster {
			b.dayMaster += tenths
		} else {
			value := tenGod(dayStem, stem)
			god = &value
			b.gods[value.Code] += tenths
		}
		b.score += tenths
		b.contributions = append(b.contributions, model.ElementContribution{
			Pillar: pillar, Source: source, HeavenlyStem: utils.CalculateHeavenlyStem(stem),
			EarthlyBranch: branch, Element: elementTerms[element], TenGod: god,
			Score: float64(tenths) / 10,
		})
	}
	for i, p := range pillars {
		add(pillarKeys[i], "exposed_stem", p.stem.Value, nil, 10, i == 2)
		branch := p.branch
		hidden := hiddenStemsOf(branch.Value)
		for j, stem := range hidden {
			add(pillarKeys[i], "hidden_stem", stem, &branch, hiddenWeights[len(hidden)][j], false)
		}
	}
	result := model.WeightedElements{
		Scheme: "hidden_stems_v1", LevelScheme: "score_levels_v1", Total: 8,
		Elements: make(map[string]model.WeightedElement, len(buckets)),
	}
	for i, b := range buckets {
		relation := (i - stemElement(dayStem) + 5) % 5
		terms := []model.LocalizedTerm{tenGodTerms[relation][0], tenGodTerms[relation][1]}
		sort.Slice(terms, func(a, b int) bool { return terms[a].Code < terms[b].Code })
		gods := make([]model.WeightedTenGod, 0, 2)
		for _, term := range terms {
			gods = append(gods, model.WeightedTenGod{TenGod: term, Score: float64(b.gods[term.Code]) / 10})
		}
		result.Elements[elementKeys[i]] = model.WeightedElement{
			Score: float64(b.score) / 10, Level: weightedLevel(b.score),
			DayMaster: float64(b.dayMaster) / 10, TenGods: gods,
			Contributions: b.contributions,
		}
	}
	return result
}
