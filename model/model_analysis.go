package model

// LocalizedTerm is a named concept with a stable code and its English and
// Vietnamese labels. Code is an ASCII identifier such as "ten_god.friend";
// clients should switch on Code, never on a label.
type LocalizedTerm struct {
	Code string `json:"code"`
	EN   string `json:"en"`
	VI   string `json:"vi"`
	// Nature classifies a star as NatureAuspicious, NatureInauspicious or
	// NatureMixed. It is set on stars only and omitted from every other term.
	// Schools differ on what is auspicious, and many stars cut both ways, so
	// treat it as a hint for display, not as a judgement of the chart.
	Nature string `json:"nature,omitempty"`
}

// Values of LocalizedTerm.Nature.
const (
	// NatureAuspicious marks a star that is generally read as favourable (cát).
	NatureAuspicious = "auspicious"
	// NatureInauspicious marks a star that is generally read as unfavourable (hung).
	NatureInauspicious = "inauspicious"
	// NatureMixed marks a star whose effect depends on the rest of the chart (tùy cục).
	NatureMixed = "mixed"
)

// PillarAnalysis holds the derived facts of one natal pillar. The Day Master
// pillar has IsDayMaster set and a nil TenGod.
type PillarAnalysis struct {
	IsDayMaster bool           `json:"is_day_master"`
	TenGod      *LocalizedTerm `json:"ten_god"`
	StemElement LocalizedTerm  `json:"stem_element"`
	// BranchElement is the element of the earthly branch.
	BranchElement LocalizedTerm `json:"branch_element"`
	Nayin         LocalizedTerm `json:"nayin"`
	// StemStageAtOwnBranch and StemStageAtMonthBranch are the life stages of
	// the pillar's own stem at its own branch and at the month branch.
	StemStageAtOwnBranch   LocalizedTerm `json:"stem_stage_at_own_branch"`
	StemStageAtMonthBranch LocalizedTerm `json:"stem_stage_at_month_branch"`
	// DayMasterStage is the life stage of the Day Master at this pillar's
	// branch; it is the localized form of the pillar's LifeCycle.
	DayMasterStage LocalizedTerm `json:"day_master_stage"`
	// HiddenStems lists the stems hidden in the branch, main qi first.
	HiddenStems []HiddenStem `json:"hidden_stems"`
	// IsVoid is true when the branch is one of the Day pillar's void branches.
	IsVoid bool `json:"is_void"`
	// HeavenEarthClash is true when this pillar's stem overcomes the stem of
	// another natal pillar with the same polarity and the two branches are opposite.
	HeavenEarthClash bool `json:"heaven_earth_clash"`
	// Stars holds the verified stars of this pillar in registry order, never nil.
	Stars []LocalizedTerm `json:"stars"`
}

// HiddenStem is a stem hidden in an earthly branch with its Ten God and its
// life stage at the month branch.
type HiddenStem struct {
	HeavenlyStem       HeavenlyStem  `json:"heavenly_stem"`
	TenGod             LocalizedTerm `json:"ten_god"`
	StageAtMonthBranch LocalizedTerm `json:"stage_at_month_branch"`
}

// AuxiliaryPillar is a stem and branch pair derived from the chart, such as
// the conception pillar.
type AuxiliaryPillar struct {
	HeavenlyStem  HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch EarthlyBranch `json:"earthly_branch"`
	GanZhi        GanZhi        `json:"gan_zhi"`
	Nayin         LocalizedTerm `json:"nayin"`
}

// DayMaster describes the stem of the Day pillar.
type DayMaster struct {
	HeavenlyStem HeavenlyStem  `json:"heavenly_stem"`
	Element      LocalizedTerm `json:"element"`
	Polarity     LocalizedTerm `json:"polarity"`
}

// ElementCount counts occurrences of each of the five elements.
type ElementCount struct {
	Wood  int `json:"wood"`
	Fire  int `json:"fire"`
	Earth int `json:"earth"`
	Metal int `json:"metal"`
	Water int `json:"water"`
}

// ElementCounts holds unweighted element counts of the exposed stems, the
// branches and the hidden stems, kept separate.
type ElementCounts struct {
	Stems       ElementCount `json:"stems"`
	Branches    ElementCount `json:"branches"`
	HiddenStems ElementCount `json:"hidden_stems"`
}

// ChartAnalysis holds chart-level derived facts.
type ChartAnalysis struct {
	DayMaster DayMaster `json:"day_master"`
	// VoidBranches are the two branches missing from the Day pillar's ten-pillar cycle.
	VoidBranches []EarthlyBranch `json:"void_branches"`
	// ThaiNguyen is the conception pillar (Tai Yuan); ThaiTuc is the
	// gestation pillar (Tai Xi); LifePalace is the life palace (Ming Gong).
	ThaiNguyen    AuxiliaryPillar `json:"thai_nguyen"`
	ThaiTuc       AuxiliaryPillar `json:"thai_tuc"`
	LifePalace    AuxiliaryPillar `json:"life_palace"`
	ElementCounts ElementCounts   `json:"element_counts"`
	// DayMasterStrength and UsefulGod are reserved and always null.
	DayMasterStrength any `json:"day_master_strength"`
	UsefulGod         any `json:"useful_god"`
	// YearReference is display-only: the Lichun of the birth year and the year of the
	// lunar calendar. Nothing else in the chart is calculated from it.
	YearReference *YearReference `json:"year_reference"`
}

// LuckPillarAnalysis holds the derived facts of one luck pillar. Ages are
// nominal: the year of the year pillar counts as age 1.
type LuckPillarAnalysis struct {
	AgeStart             int           `json:"age_start"`
	AgeEnd               int           `json:"age_end"`
	TenGod               LocalizedTerm `json:"ten_god"`
	Nayin                LocalizedTerm `json:"nayin"`
	StemStageAtOwnBranch LocalizedTerm `json:"stem_stage_at_own_branch"`
	// HiddenStems lists the stems hidden in the branch, main qi first, with
	// their life stage at the luck pillar's own branch.
	HiddenStems []LuckHiddenStem `json:"hidden_stems"`
	// HeavenEarthClash is true when the luck pillar clashes with a natal pillar.
	HeavenEarthClash bool `json:"heaven_earth_clash"`
}

// LuckHiddenStem is a stem hidden in a luck pillar's branch.
type LuckHiddenStem struct {
	HeavenlyStem     HeavenlyStem  `json:"heavenly_stem"`
	TenGod           LocalizedTerm `json:"ten_god"`
	StageAtOwnBranch LocalizedTerm `json:"stage_at_own_branch"`
}

// AnnualPillar is the pillar of one Gregorian year with its relation to the chart.
type AnnualPillar struct {
	Year int `json:"year"`
	// NominalAge counts the year of the year pillar as 1.
	NominalAge           int           `json:"nominal_age"`
	HeavenlyStem         HeavenlyStem  `json:"heavenly_stem"`
	EarthlyBranch        EarthlyBranch `json:"earthly_branch"`
	GanZhi               GanZhi        `json:"gan_zhi"`
	Nayin                LocalizedTerm `json:"nayin"`
	TenGod               LocalizedTerm `json:"ten_god"`
	StemStageAtOwnBranch LocalizedTerm `json:"stem_stage_at_own_branch"`
	// LuckPillarNumber is the number of the luck pillar whose years contain
	// this year, or nil when no luck pillar does.
	LuckPillarNumber *int `json:"luck_pillar_number"`
	// HeavenEarthClash is true when the year clashes with a natal pillar.
	HeavenEarthClash bool `json:"heaven_earth_clash"`
}

// AnnualPillars is the result of GetAnnualPillars.
type AnnualPillars struct {
	FromYear      int            `json:"from_year"`
	Count         int            `json:"count"`
	AnnualPillars []AnnualPillar `json:"annual_pillars"`
}
