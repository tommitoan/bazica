package bazica

type Element struct {
	Wood  Wood  `json:"wood"`
	Fire  Fire  `json:"fire"`
	Earth Earth `json:"earth"`
	Metal Metal `json:"metal"`
	Water Water `json:"water"`
}

type Wood struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
type Fire struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
type Earth struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
type Metal struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
type Water struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type SolarTermYear struct {
	MinorCold          string `json:"minor_cold"`
	MajorCold          string `json:"major_cold"`
	StartOfSpring      string `json:"start_of_spring"`
	SpringShowers      string `json:"spring_showers"`
	AwakeningOfInsects string `json:"awakening_of_insects"`
	SpringEquinox      string `json:"spring_equinox"`
	PureBrightness     string `json:"pure_brightness"`
	GrainRain          string `json:"grain_rain"`
	StartOfSummer      string `json:"start_of_summer"`
	GrainBuds          string `json:"grain_buds"`
	GrainInEar         string `json:"grain_in_ear"`
	SummerSolstice     string `json:"summer_solstice"`
	MinorHeat          string `json:"minor_heat"`
	MajorHeat          string `json:"major_heat"`
	StartOfAutumn      string `json:"start_of_autumn"`
	EndOfHeat          string `json:"end_of_heat"`
	WhiteDew           string `json:"white_dew"`
	AutumnEquinox      string `json:"autumn_equinox"`
	ColdDew            string `json:"cold_dew"`
	Frost              string `json:"frost"`
	StartOfWinter      string `json:"start_of_winter"`
	MinorSnow          string `json:"minor_snow"`
	MajorSnow          string `json:"major_snow"`
	WinterSolstice     string `json:"winter_solstice"`
}

type CombinedData struct {
	Year string        `json:"year"` // Year extracted from filename
	Data SolarTermYear `json:"data"` // Data for the specific year
}

var combinedData map[string]CombinedData

type LunarNewYearData struct {
	LunarNewYearDates map[string]string `json:"lunarNewYearDates"`
}
