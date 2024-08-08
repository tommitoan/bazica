package ultis

import "github.com/tommitoan/bazica/model"

// CalculateHeavenlyStem convert value to Heavenly Stem
func CalculateHeavenlyStem(num int) model.HeavenlyStem {
	var heavenlyStem model.HeavenlyStem

	// calculate heavenly stem (by the last digit of the year)
	switch num {
	case model.YangMetalValue:
		heavenlyStem.Name = model.YangMetalName
		heavenlyStem.Value = model.YangMetalValue
	case model.YinMetalValue:
		heavenlyStem.Name = model.YinMetalName
		heavenlyStem.Value = model.YinMetalValue
	case model.YangWaterValue:
		heavenlyStem.Name = model.YangWaterName
		heavenlyStem.Value = model.YangWaterValue
	case model.YinWaterValue:
		heavenlyStem.Name = model.YinWaterName
		heavenlyStem.Value = model.YinWaterValue
	case model.YangWoodValue:
		heavenlyStem.Name = model.YangWoodName
		heavenlyStem.Value = model.YangWoodValue
	case model.YinWoodValue:
		heavenlyStem.Name = model.YinWoodName
		heavenlyStem.Value = model.YinWoodValue
	case model.YangFireValue:
		heavenlyStem.Name = model.YangFireName
		heavenlyStem.Value = model.YangFireValue
	case model.YinFireValue:
		heavenlyStem.Name = model.YinFireName
		heavenlyStem.Value = model.YinFireValue
	case model.YangEarthValue:
		heavenlyStem.Name = model.YangEarthName
		heavenlyStem.Value = model.YangEarthValue
	case model.YinEarthValue:
		heavenlyStem.Name = model.YinEarthName
		heavenlyStem.Value = model.YinEarthValue
	}
	return heavenlyStem
}

// CalculateEarthlyBranch convert value to Earthly Branch
func CalculateEarthlyBranch(num int) model.EarthlyBranch {
	var earthlyBranch model.EarthlyBranch
	switch num {
	case model.TigerValue:
		earthlyBranch.Name = model.Tiger
		earthlyBranch.Value = model.TigerValue
	case model.RabbitValue:
		earthlyBranch.Name = model.Rabbit
		earthlyBranch.Value = model.RabbitValue
	case model.DragonValue:
		earthlyBranch.Name = model.Dragon
		earthlyBranch.Value = model.DragonValue
	case model.SnakeValue:
		earthlyBranch.Name = model.Snake
		earthlyBranch.Value = model.SnakeValue
	case model.HorseValue:
		earthlyBranch.Name = model.Horse
		earthlyBranch.Value = model.HorseValue
	case model.GoatValue:
		earthlyBranch.Name = model.Goat
		earthlyBranch.Value = model.GoatValue
	case model.MonkeyValue:
		earthlyBranch.Name = model.Monkey
		earthlyBranch.Value = model.MonkeyValue
	case model.RoosterValue:
		earthlyBranch.Name = model.Rooster
		earthlyBranch.Value = model.RoosterValue
	case model.DogValue:
		earthlyBranch.Name = model.Dog
		earthlyBranch.Value = model.DogValue
	case model.PigValue:
		earthlyBranch.Name = model.Pig
		earthlyBranch.Value = model.PigValue
	case model.RatValue:
		earthlyBranch.Name = model.Rat
		earthlyBranch.Value = model.RatValue
	case model.OxValue:
		earthlyBranch.Name = model.Ox
		earthlyBranch.Value = model.OxValue
	}
	return earthlyBranch
}

// ConvertTermToBranch convert solar term to earthly branch
// (VN) Chuyển đổi Tiết khí của năm sang Địa chi (của tháng)
func ConvertTermToBranch(termName string) model.EarthlyBranch {
	var earthlyBranch model.EarthlyBranch
	switch termName {
	case model.StartOfSpring, model.SpringShowers:
		earthlyBranch.Name = model.Tiger
		earthlyBranch.Value = model.TigerValue
	case model.AwakeningOfInsects, model.SpringEquinox:
		earthlyBranch.Name = model.Rabbit
		earthlyBranch.Value = model.RabbitValue
	case model.PureBrightness, model.GrainRain:
		earthlyBranch.Name = model.Dragon
		earthlyBranch.Value = model.DragonValue
	case model.StartOfSummer, model.GrainBuds:
		earthlyBranch.Name = model.Snake
		earthlyBranch.Value = model.SnakeValue
	case model.GrainInEar, model.SummerSolstice:
		earthlyBranch.Name = model.Horse
		earthlyBranch.Value = model.HorseValue
	case model.MinorHeat, model.MajorHeat:
		earthlyBranch.Name = model.Goat
		earthlyBranch.Value = model.GoatValue
	case model.StartOfAutumn, model.EndOfHeat:
		earthlyBranch.Name = model.Monkey
		earthlyBranch.Value = model.MonkeyValue
	case model.WhiteDew, model.AutumnEquinox:
		earthlyBranch.Name = model.Rooster
		earthlyBranch.Value = model.RoosterValue
	case model.ColdDew, model.Frost:
		earthlyBranch.Name = model.Dog
		earthlyBranch.Value = model.DogValue
	case model.StartOfWinter, model.MinorSnow:
		earthlyBranch.Name = model.Pig
		earthlyBranch.Value = model.PigValue
	case model.MajorSnow, model.WinterSolstice:
		earthlyBranch.Name = model.Rat
		earthlyBranch.Value = model.RatValue
	case model.MinorCold, model.MajorCold:
		earthlyBranch.Name = model.Ox
		earthlyBranch.Value = model.OxValue
	}
	return earthlyBranch
}

func ConvertValueToLifeCycle(value int) string {
	var lifeCycle string
	switch value {
	case 1:
		lifeCycle = model.LC1
	case 2:
		lifeCycle = model.LC2
	case 3:
		lifeCycle = model.LC3
	case 4:
		lifeCycle = model.LC4
	case 5:
		lifeCycle = model.LC5
	case 6:
		lifeCycle = model.LC6
	case 7:
		lifeCycle = model.LC7
	case 8:
		lifeCycle = model.LC8
	case 9:
		lifeCycle = model.LC9
	case 10:
		lifeCycle = model.LC10
	case 11:
		lifeCycle = model.LC11
	case 12:
		lifeCycle = model.LC12
	}
	return lifeCycle
}

// GetStemRuleByFiveTigers use Five-tigers Seek to get the Stem of Tiger month base on Year's Stem
// https://www.chinesefortunecalendar.com/Five-Tigers-Year-Month-Table.htm
// (VN) Dùng Ngũ hổ độn (ngũ dần) để tính Can Tháng dựa trên Can Năm
func GetStemRuleByFiveTigers(yearValue int) int {
	var stemValueOfFirstMonth int
	switch yearValue {
	case 1, 6:
		stemValueOfFirstMonth = 3
	case 2, 7:
		stemValueOfFirstMonth = 5
	case 3, 8:
		stemValueOfFirstMonth = 7
	case 4, 9:
		stemValueOfFirstMonth = 9
	case 5, 10:
		stemValueOfFirstMonth = 1
	}
	return stemValueOfFirstMonth
}

// GetStemRuleByFiveRats use Five-rats Seek to get the Stem of Rat hour base on Day's Stem
// https://www.chinesefortunecalendar.com/Five-Rats-Day-Hour-Table.htm
// (VN) Dùng Ngũ tí độn để tính Can Giờ dựa trên Can Ngày
func GetStemRuleByFiveRats(dayValue int) int {
	var stemValueOfRatHour int
	switch dayValue {
	case 1, 6:
		stemValueOfRatHour = 1
	case 2, 7:
		stemValueOfRatHour = 3
	case 3, 8:
		stemValueOfRatHour = 5
	case 4, 9:
		stemValueOfRatHour = 7
	case 5, 10:
		stemValueOfRatHour = 9
	}
	return stemValueOfRatHour
}

// GetLifeCycleRule gets the Life Stages value of the Tiger branch
// https://en.wikibooks.org/wiki/Ba_Zi/Life_Cycle
// (VN) Tính Vòng tràng sinh dựa trên Can Ngày và 4 Chi
func GetLifeCycleRule(dayValue int) int {
	var lifeCycleValue int
	switch dayValue {
	case 1:
		lifeCycleValue = 4
	case 2:
		lifeCycleValue = 5
	case 3, 5:
		lifeCycleValue = 1
	case 4, 6:
		lifeCycleValue = 8
	case 7:
		lifeCycleValue = 10
	case 8:
		lifeCycleValue = 11
	case 9:
		lifeCycleValue = 7
	case 10:
		lifeCycleValue = 2
	}
	return lifeCycleValue
}
