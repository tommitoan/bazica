package bazica

func GetMonthPillar(path, input string, yearPillar *YearPillar) (*MonthPillar, error) {
	var monthPillar MonthPillar
	// Detect solar term
	termName, err := DetectSolarTerm(path, input)
	if err != nil {
		return nil, err
	}

	// Get earthly branch
	earthBranch := convertTermToBranch(termName)
	monthPillar.EarthlyBranch = earthBranch

	// Get heavenly stem
	valueOfFirstMonth := calculateHeavenlyStemRule(yearPillar.HeavenlyStem.Value)
	valueToCal := (valueOfFirstMonth - 1) + monthPillar.EarthlyBranch.Value
	if valueToCal > 12 {
		valueToCal = valueToCal - 12
	}
	heavenlyStem := CalculateHeavenlyStem(valueToCal)
	monthPillar.HeavenlyStem = heavenlyStem

	solarDate, err := GetDateFormat(input)
	if err != nil {
		return nil, err
	}
	monthPillar.Month = int(solarDate.Month())

	return &monthPillar, nil
}

// Chuyển đổi Tiết khí của năm sang Địa chi (của tháng)
// convertTermToBranch convert solar term to earthly branche
func convertTermToBranch(termName string) EarthlyBranch {
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

func calculateHeavenlyStemRule(yearValue int) int {
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

func getStem(monthValue int) HeavenlyStem {
	if monthValue > 12 {
		monthValue = monthValue - 12
	}
	var heavenlyStem HeavenlyStem
	switch monthValue {
	case YangWoodValue:
		heavenlyStem.Value = YangWoodValue
		heavenlyStem.Name = YangWoodName
	case YinWoodValue:
		heavenlyStem.Value = YinWoodValue
		heavenlyStem.Name = YinWoodName
	case YangFireValue:
		heavenlyStem.Value = YangFireValue
		heavenlyStem.Name = YangFireName
	case YinFireValue:
		heavenlyStem.Value = YinFireValue
		heavenlyStem.Name = YinFireName

	}
	return heavenlyStem
}
