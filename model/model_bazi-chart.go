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
	LuckPillars []*LuckPillar `json:"luck_pillars"`
}

type LuckPillar struct {
	Number        int           `json:"number"`
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
