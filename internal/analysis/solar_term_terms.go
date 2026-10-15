package analysis

import "github.com/tommitoan/bazica/v2/model"

// solarTermTerms gives each calendar term a stable code and bilingual display names.
var solarTermTerms = map[string]model.LocalizedTerm{
	model.MinorCold:          {Code: "solar_term.minor_cold", EN: "Minor Cold", VI: "Tiểu Hàn"},
	model.MajorCold:          {Code: "solar_term.major_cold", EN: "Major Cold", VI: "Đại Hàn"},
	model.StartOfSpring:      {Code: "solar_term.start_of_spring", EN: "Start of Spring", VI: "Lập Xuân"},
	model.SpringShowers:      {Code: "solar_term.spring_showers", EN: "Spring Showers", VI: "Vũ Thủy"},
	model.AwakeningOfInsects: {Code: "solar_term.awakening_of_insects", EN: "Awakening of Insects", VI: "Kinh Trập"},
	model.SpringEquinox:      {Code: "solar_term.spring_equinox", EN: "Spring Equinox", VI: "Xuân Phân"},
	model.PureBrightness:     {Code: "solar_term.pure_brightness", EN: "Pure Brightness", VI: "Thanh Minh"},
	model.GrainRain:          {Code: "solar_term.grain_rain", EN: "Grain Rain", VI: "Cốc Vũ"},
	model.StartOfSummer:      {Code: "solar_term.start_of_summer", EN: "Start of Summer", VI: "Lập Hạ"},
	model.GrainBuds:          {Code: "solar_term.grain_buds", EN: "Grain Buds", VI: "Tiểu Mãn"},
	model.GrainInEar:         {Code: "solar_term.grain_in_ear", EN: "Grain in Ear", VI: "Mang Chủng"},
	model.SummerSolstice:     {Code: "solar_term.summer_solstice", EN: "Summer Solstice", VI: "Hạ Chí"},
	model.MinorHeat:          {Code: "solar_term.minor_heat", EN: "Minor Heat", VI: "Tiểu Thử"},
	model.MajorHeat:          {Code: "solar_term.major_heat", EN: "Major Heat", VI: "Đại Thử"},
	model.StartOfAutumn:      {Code: "solar_term.start_of_autumn", EN: "Start of Autumn", VI: "Lập Thu"},
	model.EndOfHeat:          {Code: "solar_term.end_of_heat", EN: "End of Heat", VI: "Xử Thử"},
	model.WhiteDew:           {Code: "solar_term.white_dew", EN: "White Dew", VI: "Bạch Lộ"},
	model.AutumnEquinox:      {Code: "solar_term.autumn_equinox", EN: "Autumn Equinox", VI: "Thu Phân"},
	model.ColdDew:            {Code: "solar_term.cold_dew", EN: "Cold Dew", VI: "Hàn Lộ"},
	model.Frost:              {Code: "solar_term.frost", EN: "Frost Descent", VI: "Sương Giáng"},
	model.StartOfWinter:      {Code: "solar_term.start_of_winter", EN: "Start of Winter", VI: "Lập Đông"},
	model.MinorSnow:          {Code: "solar_term.minor_snow", EN: "Minor Snow", VI: "Tiểu Tuyết"},
	model.MajorSnow:          {Code: "solar_term.major_snow", EN: "Major Snow", VI: "Đại Tuyết"},
	model.WinterSolstice:     {Code: "solar_term.winter_solstice", EN: "Winter Solstice", VI: "Đông Chí"},
}
