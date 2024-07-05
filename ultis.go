package bazica

// CalculateHeavenlyStem convert value to Heavenly Stem
func CalculateHeavenlyStem(num int) HeavenlyStem {
	var heavenlyStem HeavenlyStem

	// calculate heavenly stem (by the last digit of the year)
	switch num {
	case YangMetalValue:
		heavenlyStem.Name = YangMetalName
		heavenlyStem.Value = YangMetalValue
	case YinMetalValue:
		heavenlyStem.Name = YinMetalName
		heavenlyStem.Value = YinMetalValue
	case YangWaterValue:
		heavenlyStem.Name = YangWaterName
		heavenlyStem.Value = YangWaterValue
	case YinWaterValue:
		heavenlyStem.Name = YinWaterName
		heavenlyStem.Value = YinWaterValue
	case YangWoodValue:
		heavenlyStem.Name = YangWoodName
		heavenlyStem.Value = YangWoodValue
	case YinWoodValue:
		heavenlyStem.Name = YinWoodName
		heavenlyStem.Value = YinWoodValue
	case YangFireValue:
		heavenlyStem.Name = YangFireName
		heavenlyStem.Value = YangFireValue
	case YinFireValue:
		heavenlyStem.Name = YinFireName
		heavenlyStem.Value = YinFireValue
	case YangEarthValue:
		heavenlyStem.Name = YangEarthName
		heavenlyStem.Value = YangEarthValue
	case YinEarthValue:
		heavenlyStem.Name = YinEarthName
		heavenlyStem.Value = YinEarthValue
	}
	return heavenlyStem
}

// CalculateEarthlyBranch convert value to Earthly Branch
func CalculateEarthlyBranch(num int) EarthlyBranch {
	var earthlyBranch EarthlyBranch
	switch num {
	case TigerValue:
		earthlyBranch.Name = Tiger
		earthlyBranch.Value = TigerValue
	case RabbitValue:
		earthlyBranch.Name = Rabbit
		earthlyBranch.Value = RabbitValue
	case DragonValue:
		earthlyBranch.Name = Dragon
		earthlyBranch.Value = DragonValue
	case SnakeValue:
		earthlyBranch.Name = Snake
		earthlyBranch.Value = SnakeValue
	case HorseValue:
		earthlyBranch.Name = Horse
		earthlyBranch.Value = HorseValue
	case GoatValue:
		earthlyBranch.Name = Goat
		earthlyBranch.Value = GoatValue
	case MonkeyValue:
		earthlyBranch.Name = Monkey
		earthlyBranch.Value = MonkeyValue
	case RoosterValue:
		earthlyBranch.Name = Rooster
		earthlyBranch.Value = RoosterValue
	case DogValue:
		earthlyBranch.Name = Dog
		earthlyBranch.Value = DogValue
	case PigValue:
		earthlyBranch.Name = Pig
		earthlyBranch.Value = PigValue
	case RatValue:
		earthlyBranch.Name = Rat
		earthlyBranch.Value = RatValue
	case OxValue:
		earthlyBranch.Name = Ox
		earthlyBranch.Value = OxValue
	}
	return earthlyBranch
}

// ConvertTermToBranch convert solar term to earthly branch
// (VN) Chuyển đổi Tiết khí của năm sang Địa chi (của tháng)
func ConvertTermToBranch(termName string) EarthlyBranch {
	var earthlyBranch EarthlyBranch
	switch termName {
	case StartOfSpring, SpringShowers:
		earthlyBranch.Name = Tiger
		earthlyBranch.Value = TigerValue
	case AwakeningOfInsects, SpringEquinox:
		earthlyBranch.Name = Rabbit
		earthlyBranch.Value = RabbitValue
	case PureBrightness, GrainRain:
		earthlyBranch.Name = Dragon
		earthlyBranch.Value = DragonValue
	case StartOfSummer, GrainBuds:
		earthlyBranch.Name = Snake
		earthlyBranch.Value = SnakeValue
	case GrainInEar, SummerSolstice:
		earthlyBranch.Name = Horse
		earthlyBranch.Value = HorseValue
	case MinorHeat, MajorHeat:
		earthlyBranch.Name = Goat
		earthlyBranch.Value = GoatValue
	case StartOfAutumn, EndOfHeat:
		earthlyBranch.Name = Monkey
		earthlyBranch.Value = MonkeyValue
	case WhiteDew, AutumnEquinox:
		earthlyBranch.Name = Rooster
		earthlyBranch.Value = RoosterValue
	case ColdDew, Frost:
		earthlyBranch.Name = Dog
		earthlyBranch.Value = DogValue
	case StartOfWinter, MinorSnow:
		earthlyBranch.Name = Pig
		earthlyBranch.Value = PigValue
	case MajorSnow, WinterSolstice:
		earthlyBranch.Name = Rat
		earthlyBranch.Value = RatValue
	case MinorCold, MajorCold:
		earthlyBranch.Name = Ox
		earthlyBranch.Value = OxValue
	}
	return earthlyBranch
}

// GetStemRuleByFireTigers use Five-tigers Seek to get the Stem of Tiger month base on Year's Stem
// https://www.chinesefortunecalendar.com/Five-Tigers-Year-Month-Table.htm
// (VN) Dùng Ngũ hổ độn (ngũ dần) để tính Can Tháng dựa trên Can Năm
func GetStemRuleByFireTigers(yearValue int) int {
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

// GetStemRuleByFireRats use Five-rats Seek to get the Stem of Rat hour base on Day's Stem
// https://www.chinesefortunecalendar.com/Five-Rats-Day-Hour-Table.htm
// (VN) Dùng Ngũ tí độn để tính Can Giờ dựa trên Can Ngày
func GetStemRuleByFireRats(dayValue int) int {
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
