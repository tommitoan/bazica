package model

// LocalizedTerm is a named concept with a stable code and its English and
// Vietnamese labels. Code is an ASCII identifier such as "ten_god.friend";
// clients should switch on Code, never on a label.
type LocalizedTerm struct {
	Code string `json:"code"`
	EN   string `json:"en"`
	VI   string `json:"vi"`
}

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
	// Stars holds verified stars only; it is empty until stars are evaluated.
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
}
