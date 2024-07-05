package model

type BaziChart struct {
	YearPillar  *YearPillar  `json:"year_pillar"`
	MonthPillar *MonthPillar `json:"month_pillar"`
	DayPillar   *DayPillar   `json:"day_pillar"`
	HourPillar  *HourPillar  `json:"hour_pillar"`
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
