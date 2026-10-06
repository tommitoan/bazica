package model

import "time"

// YearReference is display-only data about the year boundary of a chart. The
// Ba-zi year pillar changes at Lichun; many people know the year of the lunar
// calendar, which changes at the Lunar New Year, so both are reported. No pillar,
// star, luck pillar, age or yearly row is calculated from this block.
type YearReference struct {
	// Lichun is the Lichun of the civil year of the birth.
	Lichun LichunReference `json:"lichun"`
	// LunarYear is the year of the lunar calendar the birth falls in.
	LunarYear LunarYearReference `json:"lunar_year"`
	// Differs is true when the lunar-calendar year is not the Ba-zi year, which happens
	// for a birth between the Lunar New Year and Lichun (about 2 % of births).
	Differs bool `json:"differs"`
}

// LichunReference is the instant the Ba-zi year changes in the birth's civil year.
type LichunReference struct {
	// Time is the instant of Lichun, expressed in the birth's time zone.
	Time time.Time `json:"time"`
	// BirthBeforeLichun is true when the birth precedes that instant.
	BirthBeforeLichun bool `json:"birth_before_lichun"`
}

// LunarYearReference is the year of the lunar calendar, named by stem and branch.
// The branch name is the zodiac animal of that year.
type LunarYearReference struct {
	// Year is the Gregorian year in which that lunar year begins.
	Year          int           `json:"year"`
	HeavenlyStem  HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch `json:"earthly_branch"`
	GanZhi        GanZhi        `json:"gan_zhi"`
	// NewYearDate is the civil date (YYYY-MM-DD) of the Lunar New Year that opens
	// that lunar year; it is empty for the lunar year before the first table year.
	NewYearDate string `json:"new_year_date,omitempty"`
}
