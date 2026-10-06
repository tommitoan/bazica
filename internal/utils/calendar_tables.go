package utils

import "github.com/tommitoan/bazica/v2/model"

// heavenlyStems is indexed by stem value - 1.
var heavenlyStems = [...]model.HeavenlyStem{
	{Name: model.YangWoodName, Value: model.YangWoodValue, Character: model.YangWoodCharacter, Spelling: model.YangWoodSpelling},
	{Name: model.YinWoodName, Value: model.YinWoodValue, Character: model.YinWoodCharacter, Spelling: model.YinWoodSpelling},
	{Name: model.YangFireName, Value: model.YangFireValue, Character: model.YangFireCharacter, Spelling: model.YangFireSpelling},
	{Name: model.YinFireName, Value: model.YinFireValue, Character: model.YinFireCharacter, Spelling: model.YinFireSpelling},
	{Name: model.YangEarthName, Value: model.YangEarthValue, Character: model.YangEarthCharacter, Spelling: model.YangEarthSpelling},
	{Name: model.YinEarthName, Value: model.YinEarthValue, Character: model.YinEarthCharacter, Spelling: model.YinEarthSpelling},
	{Name: model.YangMetalName, Value: model.YangMetalValue, Character: model.YangMetalCharacter, Spelling: model.YangMetalSpelling},
	{Name: model.YinMetalName, Value: model.YinMetalValue, Character: model.YinMetalCharacter, Spelling: model.YinMetalSpelling},
	{Name: model.YangWaterName, Value: model.YangWaterValue, Character: model.YangWaterCharacter, Spelling: model.YangWaterSpelling},
	{Name: model.YinWaterName, Value: model.YinWaterValue, Character: model.YinWaterCharacter, Spelling: model.YinWaterSpelling},
}

// earthlyBranches is indexed by branch value - 1 (Tiger first).
var earthlyBranches = [...]model.EarthlyBranch{
	{Name: model.Tiger, Value: model.TigerValue, Character: model.TigerCharacter, Spelling: model.TigerSpelling},
	{Name: model.Rabbit, Value: model.RabbitValue, Character: model.RabbitCharacter, Spelling: model.RabbitSpelling},
	{Name: model.Dragon, Value: model.DragonValue, Character: model.DragonCharacter, Spelling: model.DragonSpelling},
	{Name: model.Snake, Value: model.SnakeValue, Character: model.SnakeCharacter, Spelling: model.SnakeSpelling},
	{Name: model.Horse, Value: model.HorseValue, Character: model.HorseCharacter, Spelling: model.HorseSpelling},
	{Name: model.Goat, Value: model.GoatValue, Character: model.GoatCharacter, Spelling: model.GoatSpelling},
	{Name: model.Monkey, Value: model.MonkeyValue, Character: model.MonkeyCharacter, Spelling: model.MonkeySpelling},
	{Name: model.Rooster, Value: model.RoosterValue, Character: model.RoosterCharacter, Spelling: model.RoosterSpelling},
	{Name: model.Dog, Value: model.DogValue, Character: model.DogCharacter, Spelling: model.DogSpelling},
	{Name: model.Pig, Value: model.PigValue, Character: model.PigCharacter, Spelling: model.PigSpelling},
	{Name: model.Rat, Value: model.RatValue, Character: model.RatCharacter, Spelling: model.RatSpelling},
	{Name: model.Ox, Value: model.OxValue, Character: model.OxCharacter, Spelling: model.OxSpelling},
}

// termBranchValue maps each solar term to the branch value of the month it falls in.
var termBranchValue = map[string]int{
	model.StartOfSpring: model.TigerValue, model.SpringShowers: model.TigerValue,
	model.AwakeningOfInsects: model.RabbitValue, model.SpringEquinox: model.RabbitValue,
	model.PureBrightness: model.DragonValue, model.GrainRain: model.DragonValue,
	model.StartOfSummer: model.SnakeValue, model.GrainBuds: model.SnakeValue,
	model.GrainInEar: model.HorseValue, model.SummerSolstice: model.HorseValue,
	model.MinorHeat: model.GoatValue, model.MajorHeat: model.GoatValue,
	model.StartOfAutumn: model.MonkeyValue, model.EndOfHeat: model.MonkeyValue,
	model.WhiteDew: model.RoosterValue, model.AutumnEquinox: model.RoosterValue,
	model.ColdDew: model.DogValue, model.Frost: model.DogValue,
	model.StartOfWinter: model.PigValue, model.MinorSnow: model.PigValue,
	model.MajorSnow: model.RatValue, model.WinterSolstice: model.RatValue,
	model.MinorCold: model.OxValue, model.MajorCold: model.OxValue,
}

// lifeCycles is indexed by life stage value - 1.
var lifeCycles = [...]string{
	model.LC1, model.LC2, model.LC3, model.LC4, model.LC5, model.LC6,
	model.LC7, model.LC8, model.LC9, model.LC10, model.LC11, model.LC12,
}

// fiveTigersStem and fiveRatsStem hold the stem of the first month (Tiger) and
// of the first hour (Rat), indexed by (stem value - 1) % 5. Each pair of
// stems, Jia/Ji, Yi/Geng and so on, shares one entry.
var (
	fiveTigersStem = [5]int{3, 5, 7, 9, 1}
	fiveRatsStem   = [5]int{1, 3, 5, 7, 9}
)

// lifeCycleOfTiger is the life stage of the Tiger branch for each Day Master,
// indexed by stem value - 1.
var lifeCycleOfTiger = [...]int{4, 5, 1, 8, 1, 8, 10, 11, 7, 2}

// CalculateHeavenlyStem converts a stem value (1-10) into a Heavenly Stem. An
// out-of-range value yields the zero value.
func CalculateHeavenlyStem(num int) model.HeavenlyStem {
	if num < 1 || num > len(heavenlyStems) {
		return model.HeavenlyStem{}
	}
	return heavenlyStems[num-1]
}

// CalculateEarthlyBranch converts a branch value (1-12) into an Earthly
// Branch. An out-of-range value yields the zero value.
func CalculateEarthlyBranch(num int) model.EarthlyBranch {
	if num < 1 || num > len(earthlyBranches) {
		return model.EarthlyBranch{}
	}
	return earthlyBranches[num-1]
}

// ConvertTermToBranch converts a solar term to the Earthly Branch of its month.
// (VN) Chuyển đổi Tiết khí của năm sang Địa chi (của tháng)
func ConvertTermToBranch(termName string) model.EarthlyBranch {
	return CalculateEarthlyBranch(termBranchValue[termName])
}

// ConvertValueToLifeCycle converts a life stage value (1-12) to its name.
func ConvertValueToLifeCycle(value int) string {
	if value < 1 || value > len(lifeCycles) {
		return ""
	}
	return lifeCycles[value-1]
}

// stemPair returns the entry of a five-way table for a stem value, or 0 when
// the value is not a valid stem.
func stemPair(table [5]int, stemValue int) int {
	if stemValue < 1 || stemValue > len(heavenlyStems) {
		return 0
	}
	return table[(stemValue-1)%5]
}

// GetStemRuleByFiveTigers use Five-tigers Seek to get the Stem of Tiger month base on Year's Stem
// https://www.chinesefortunecalendar.com/Five-Tigers-Year-Month-Table.htm
// (VN) Dùng Ngũ hổ độn (ngũ dần) để tính Can Tháng dựa trên Can Năm
func GetStemRuleByFiveTigers(yearValue int) int {
	return stemPair(fiveTigersStem, yearValue)
}

// GetStemRuleByFiveRats use Five-rats Seek to get the Stem of Rat hour base on Day's Stem
// https://www.chinesefortunecalendar.com/Five-Rats-Day-Hour-Table.htm
// (VN) Dùng Ngũ tí độn để tính Can Giờ dựa trên Can Ngày
func GetStemRuleByFiveRats(dayValue int) int {
	return stemPair(fiveRatsStem, dayValue)
}

// GetLifeCycleRule gets the Life Stages value of the Tiger branch
// https://en.wikibooks.org/wiki/Ba_Zi/Life_Cycle
// (VN) Tính Vòng tràng sinh dựa trên Can Ngày và 4 Chi
func GetLifeCycleRule(dayValue int) int {
	if dayValue < 1 || dayValue > len(lifeCycleOfTiger) {
		return 0
	}
	return lifeCycleOfTiger[dayValue-1]
}
