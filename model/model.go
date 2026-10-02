package model

type CombinedData struct {
	Year string        `json:"year"` // Year extracted from filename
	Data SolarTermYear `json:"data"` // Data for the specific year
}

type LunarNewYearData struct {
	LunarNewYearDates map[string]string `json:"lunarNewYearDates"`
}
