package bazica

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
