package ultis

import "github.com/tommitoan/bazica/model"

func stemElementScore(stemValue int, score *model.FiveElements) {
	switch stemValue {
	case model.YangWoodValue, model.YinWoodValue:
		score.Wood += 1
	case model.YangFireValue, model.YinFireValue:
		score.Fire += 1
	case model.YangEarthValue, model.YinEarthValue:
		score.Earth += 1
	case model.YangMetalValue, model.YinMetalValue:
		score.Metal += 1
	case model.YangWaterValue, model.YinWaterValue:
		score.Water += 1
	}
}

// branchFlexScore adds +1 to the primary element of an EarthlyBranch.
func branchFlexScore(branchValue int, score *model.FiveElements) {
	switch branchValue {
	case model.TigerValue:
		score.Wood += 1
	case model.RabbitValue:
		score.Wood += 1
	case model.DragonValue:
		score.Earth += 1
	case model.SnakeValue:
		score.Fire += 1
	case model.HorseValue:
		score.Fire += 1
	case model.GoatValue:
		score.Earth += 1
	case model.MonkeyValue:
		score.Metal += 1
	case model.RoosterValue:
		score.Metal += 1
	case model.DogValue:
		score.Earth += 1
	case model.PigValue:
		score.Water += 1
	case model.RatValue:
		score.Water += 1
	case model.OxValue:
		score.Earth += 1
	}
}

// branchStrictScore distributes an EarthlyBranch across its hidden stems (Tàng Can).
func branchStrictScore(branchValue int, score *model.FiveElements) {
	switch branchValue {
	case model.TigerValue: // Tiger: Yang Wood 60%, Yang Fire 30%, Yang Earth 10%
		score.Wood += 0.6
		score.Fire += 0.3
		score.Earth += 0.1
	case model.MonkeyValue: // Monkey: Yang Metal 60%, Yang Water 30%, Yang Earth 10%
		score.Metal += 0.6
		score.Water += 0.3
		score.Earth += 0.1
	case model.SnakeValue: // Snake: Yang Fire 60%, Yang Metal 30%, Yang Earth 10%
		score.Fire += 0.6
		score.Metal += 0.3
		score.Earth += 0.1
	case model.PigValue: // Pig: Yang Water 70%, Yang Wood 30%
		score.Water += 0.7
		score.Wood += 0.3
	case model.RatValue: // Rat: Yin Water 100%
		score.Water += 1
	case model.HorseValue: // Horse: Yin Fire 70%, Yin Earth 30%
		score.Fire += 0.7
		score.Earth += 0.3
	case model.RabbitValue: // Rabbit: Yin Wood 100%
		score.Wood += 1
	case model.RoosterValue: // Rooster: Yin Metal 100%
		score.Metal += 1
	case model.DragonValue: // Dragon: Yang Earth 60%, Yin Wood 30%, Yin Water 10%
		score.Earth += 0.6
		score.Wood += 0.3
		score.Water += 0.1
	case model.DogValue: // Dog: Yang Earth 60%, Yin Metal 30%, Yin Fire 10%
		score.Earth += 0.6
		score.Metal += 0.3
		score.Fire += 0.1
	case model.OxValue: // Ox: Yin Earth 60%, Yin Water 30%, Yin Metal 10%
		score.Earth += 0.6
		score.Water += 0.3
		score.Metal += 0.1
	case model.GoatValue: // Goat: Yin Earth 60%, Yin Fire 30%, Yin Wood 10%
		score.Earth += 0.6
		score.Fire += 0.3
		score.Wood += 0.1
	}
}

func GetElementScoreFromFourPillar(pillars *model.FourPillars) *model.ElementScore {
	strict := &model.FiveElements{}
	flex := &model.FiveElements{}
	type pillarPair struct{ stem, branch int }
	pairs := []pillarPair{
		{pillars.YearPillar.HeavenlyStem.Value, pillars.YearPillar.EarthlyBranch.Value},
		{pillars.MonthPillar.HeavenlyStem.Value, pillars.MonthPillar.EarthlyBranch.Value},
		{pillars.DayPillar.HeavenlyStem.Value, pillars.DayPillar.EarthlyBranch.Value},
		{pillars.HourPillar.HeavenlyStem.Value, pillars.HourPillar.EarthlyBranch.Value},
	}
	for _, p := range pairs {
		stemElementScore(p.stem, strict)
		stemElementScore(p.stem, flex)
		branchStrictScore(p.branch, strict)
		branchFlexScore(p.branch, flex)
	}
	return &model.ElementScore{
		StrictScore: strict,
		FlexScore:   flex,
	}
}
