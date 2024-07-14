package model

type BaziChart struct {
	PersonalInfo *PersonalInfo `json:"personalInfo,omitempty"`
	FourPillar   *FourPillars  `json:"four_pillars"`
	LuckPillars  *LuckPillars  `json:"luck_pillars"`
}

type PersonalInfo struct {
	Name   string `json:"name"`
	Gender int    `json:"gender"`
}

type FourPillars struct {
	YearPillar  *YearPillar  `json:"year_pillar"`
	MonthPillar *MonthPillar `json:"month_pillar"`
	DayPillar   *DayPillar   `json:"day_pillar"`
	HourPillar  *HourPillar  `json:"hour_pillar"`
}

type LuckPillars struct {
	LuckPillar0  *LuckPillar `json:"luck_pillar_0"`
	LuckPillar1  *LuckPillar `json:"luck_pillar_1"`
	LuckPillar2  *LuckPillar `json:"luck_pillar_2"`
	LuckPillar3  *LuckPillar `json:"luck_pillar_3"`
	LuckPillar4  *LuckPillar `json:"luck_pillar_4"`
	LuckPillar5  *LuckPillar `json:"luck_pillar_5"`
	LuckPillar6  *LuckPillar `json:"luck_pillar_6"`
	LuckPillar7  *LuckPillar `json:"luck_pillar_7"`
	LuckPillar8  *LuckPillar `json:"luck_pillar_8"`
	LuckPillar9  *LuckPillar `json:"luck_pillar_9"`
	LuckPillar10 *LuckPillar `json:"luck_pillar_10"`
	LuckPillar11 *LuckPillar `json:"luck_pillar_11"`
	LuckPillar12 *LuckPillar `json:"luck_pillar_12"`
}

type LuckPillar struct {
	HeavenlyStem  HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch `json:"earthly_branch"`
	Year          int           `json:"year"`
	Age           int           `json:"age"`
}

type YearPillar struct {
	HeavenlyStem  HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch `json:"earthly_branch"`
	Year          int           `json:"year"`
}

type MonthPillar struct {
	HeavenlyStem  HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch `json:"earthly_branch"`
	Month         int           `json:"month"`
}

type DayPillar struct {
	HeavenlyStem  HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch `json:"earthly_branch"`
	Day           int           `json:"day"`
}

type HourPillar struct {
	HeavenlyStem  HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch `json:"earthly_branch"`
	Hour          TimeOfDay     `json:"hour"`
}

type TimeOfDay struct {
	Hour   int
	Minute int
}

type HeavenlyStem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type EarthlyBranch struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
