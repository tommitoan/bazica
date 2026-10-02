package model

import "time"

// BaziChart is the result of GetBaziChart.
type BaziChart struct {
	PersonalInfo *PersonalInfo `json:"personalInfo,omitempty"`
	FourPillar   *FourPillars  `json:"four_pillars"`
	LuckPillars  *LuckPillars  `json:"luck_pillars"`
	// Analysis holds chart-level derived facts such as the Day Master, void branches and element counts.
	Analysis *ChartAnalysis `json:"analysis"`
}

// PersonalInfo is optional metadata a caller can attach to a chart; GetBaziChart does not fill it in.
type PersonalInfo struct {
	Name   string `json:"name"`
	Gender int    `json:"gender"`
}

// FourPillars holds the year, month, day and hour pillars of a birth time.
type FourPillars struct {
	YearPillar  *YearPillar  `json:"year_pillar"`
	MonthPillar *MonthPillar `json:"month_pillar"`
	DayPillar   *DayPillar   `json:"day_pillar"`
	HourPillar  *HourPillar  `json:"hour_pillar"`
}

// LuckPillars is the sequence of ten-year luck periods (Da Yun).
type LuckPillars struct {
	LuckPillars []*LuckPillar `json:"luck_pillars"`
}

// LuckPillar is one ten-year luck period. Number 0 is the birth month pillar,
// which applies until the first period starts at Time.
type LuckPillar struct {
	Number        int           `json:"number"`
	HeavenlyStem  HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch `json:"earthly_branch"`
	GanZhi        GanZhi        `json:"gan_zhi"`
	YearStart     int           `json:"year_start"`
	YearEnd       int           `json:"year_end"`
	Time          time.Time     `json:"time"`
	// Analysis holds the derived facts of the luck pillar.
	Analysis *LuckPillarAnalysis `json:"analysis"`
}

// YearPillar is the pillar of the lunar year the birth belongs to; Year is the Gregorian year.
type YearPillar struct {
	HeavenlyStem  HeavenlyStem    `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch   `json:"earthly_branch"`
	GanZhi        GanZhi          `json:"gan_zhi"`
	Year          int             `json:"year"`
	LifeCycle     string          `json:"life_cycle"`
	Analysis      *PillarAnalysis `json:"analysis"`
}

// MonthPillar is the pillar of the solar-term month; Month is the Gregorian month.
type MonthPillar struct {
	HeavenlyStem  HeavenlyStem    `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch   `json:"earthly_branch"`
	GanZhi        GanZhi          `json:"gan_zhi"`
	Month         int             `json:"month"`
	LifeCycle     string          `json:"life_cycle"`
	Analysis      *PillarAnalysis `json:"analysis"`
}

// DayPillar is the pillar of the day; the day changes at 23:00.
type DayPillar struct {
	HeavenlyStem  HeavenlyStem    `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch   `json:"earthly_branch"`
	GanZhi        GanZhi          `json:"gan_zhi"`
	Day           int             `json:"day"`
	LifeCycle     string          `json:"life_cycle"`
	Analysis      *PillarAnalysis `json:"analysis"`
}

// HourPillar is the pillar of the two-hour period containing the birth time.
type HourPillar struct {
	HeavenlyStem  HeavenlyStem    `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch   `json:"earthly_branch"`
	GanZhi        GanZhi          `json:"gan_zhi"`
	Hour          TimeOfDay       `json:"hour"`
	LifeCycle     string          `json:"life_cycle"`
	Analysis      *PillarAnalysis `json:"analysis"`
}

// GanZhi is the Nayin (sound element) of a stem and branch pair.
type GanZhi struct {
	Name         string `json:"name"`
	ElementName  string `json:"element_name"`
	ElementValue int    `json:"element_value"`
}

// TimeOfDay is a wall-clock time.
type TimeOfDay struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

// HeavenlyStem is one of the ten stems; Value runs from 1 (Yang Wood, Jia) to 10 (Yin Water, Gui).
type HeavenlyStem struct {
	Name      string `json:"name"`
	Value     int    `json:"value"`
	Character string `json:"character"`
	Spelling  string `json:"spelling"`
}

// EarthlyBranch is one of the twelve branches; Value runs from 1 (Tiger) to 12 (Ox).
type EarthlyBranch struct {
	Name      string `json:"name"`
	Value     int    `json:"value"`
	Character string `json:"character"`
	Spelling  string `json:"spelling"`
}
