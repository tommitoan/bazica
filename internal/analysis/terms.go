package analysis

import "github.com/tommitoan/bazica/model"

// The tables below are the bilingual term registry. Codes are stable identifiers;
// the Nayin English labels equal the names returned by utils.GetGanzhi.

// tenGodTerms is indexed by relation to the Day Master (same element, generated,
// controlled, controlling, generating) and then by polarity (same, opposite).
var tenGodTerms = [5][2]model.LocalizedTerm{
	{{Code: "ten_god.friend", EN: "Friend", VI: "Tỷ Kiên"}, {Code: "ten_god.rob_wealth", EN: "Rob Wealth", VI: "Kiếp Tài"}},
	{{Code: "ten_god.eating_god", EN: "Eating God", VI: "Thực Thần"}, {Code: "ten_god.hurting_officer", EN: "Hurting Officer", VI: "Thương Quan"}},
	{{Code: "ten_god.indirect_wealth", EN: "Indirect Wealth", VI: "Thiên Tài"}, {Code: "ten_god.direct_wealth", EN: "Direct Wealth", VI: "Chính Tài"}},
	{{Code: "ten_god.seven_killings", EN: "Seven Killings", VI: "Thất Sát"}, {Code: "ten_god.direct_officer", EN: "Direct Officer", VI: "Chính Quan"}},
	{{Code: "ten_god.indirect_resource", EN: "Indirect Resource", VI: "Thiên Ấn"}, {Code: "ten_god.direct_resource", EN: "Direct Resource", VI: "Chính Ấn"}},
}

// stageTerms is indexed by life stage order, Birth first.
var stageTerms = [12]model.LocalizedTerm{
	{Code: "stage.birth", EN: "Birth", VI: "Sinh"},
	{Code: "stage.bath", EN: "Bath", VI: "Mộc Dục"},
	{Code: "stage.youth", EN: "Youth", VI: "Quan Đới"},
	{Code: "stage.thriving", EN: "Thriving", VI: "Lâm Quan"},
	{Code: "stage.prosperous", EN: "Prosperous", VI: "Đế Vượng"},
	{Code: "stage.weakening", EN: "Weakening", VI: "Suy"},
	{Code: "stage.sick", EN: "Sick", VI: "Bệnh"},
	{Code: "stage.death", EN: "Death", VI: "Tử"},
	{Code: "stage.grave", EN: "Grave", VI: "Mộ"},
	{Code: "stage.extinction", EN: "Extinction", VI: "Tuyệt"},
	{Code: "stage.conceived", EN: "Conceived", VI: "Thai"},
	{Code: "stage.nourishing", EN: "Nourishing", VI: "Dưỡng"},
}

// elementTerms is indexed Wood, Fire, Earth, Metal, Water.
var elementTerms = [5]model.LocalizedTerm{
	{Code: "element.wood", EN: "Wood", VI: "Mộc"},
	{Code: "element.fire", EN: "Fire", VI: "Hỏa"},
	{Code: "element.earth", EN: "Earth", VI: "Thổ"},
	{Code: "element.metal", EN: "Metal", VI: "Kim"},
	{Code: "element.water", EN: "Water", VI: "Thủy"},
}

// polarityTerms is indexed Yang, Yin.
var polarityTerms = [2]model.LocalizedTerm{
	{Code: "polarity.yang", EN: "Yang", VI: "Dương"},
	{Code: "polarity.yin", EN: "Yin", VI: "Âm"},
}

// nayinTerms is keyed by the English Nayin name used in GanZhi.Name.
var nayinTerms = map[string]model.LocalizedTerm{
	"Sea metal":        {Code: "nayin.sea_metal", EN: "Sea metal", VI: "Hải Trung Kim"},
	"Furnace fire":     {Code: "nayin.furnace_fire", EN: "Furnace fire", VI: "Lò Trung Hỏa"},
	"Forest wood":      {Code: "nayin.forest_wood", EN: "Forest wood", VI: "Đại Lâm Mộc"},
	"Road earth":       {Code: "nayin.road_earth", EN: "Road earth", VI: "Lộ Bàng Thổ"},
	"Sword metal":      {Code: "nayin.sword_metal", EN: "Sword metal", VI: "Kiếm Phong Kim"},
	"Volcanic fire":    {Code: "nayin.volcanic_fire", EN: "Volcanic fire", VI: "Sơn Đầu Hỏa"},
	"Cave water":       {Code: "nayin.cave_water", EN: "Cave water", VI: "Giản Hạ Thủy"},
	"Fortress earth":   {Code: "nayin.fortress_earth", EN: "Fortress earth", VI: "Thành Đầu Thổ"},
	"Wax metal":        {Code: "nayin.wax_metal", EN: "Wax metal", VI: "Bạch Lạp Kim"},
	"Willow wood":      {Code: "nayin.willow_wood", EN: "Willow wood", VI: "Dương Liễu Mộc"},
	"Stream water":     {Code: "nayin.stream_water", EN: "Stream water", VI: "Tuyền Trung Thủy"},
	"Roof tiles earth": {Code: "nayin.roof_tiles_earth", EN: "Roof tiles earth", VI: "Ốc Thượng Thổ"},
	"Lightning fire":   {Code: "nayin.lightning_fire", EN: "Lightning fire", VI: "Tích Lịch Hỏa"},
	"Conifer wood":     {Code: "nayin.conifer_wood", EN: "Conifer wood", VI: "Tùng Bách Mộc"},
	"River water":      {Code: "nayin.river_water", EN: "River water", VI: "Trường Lưu Thủy"},
	"Sand metal":       {Code: "nayin.sand_metal", EN: "Sand metal", VI: "Sa Trung Kim"},
	"Forest fire":      {Code: "nayin.forest_fire", EN: "Forest fire", VI: "Sơn Hạ Hỏa"},
	"Meadow wood":      {Code: "nayin.meadow_wood", EN: "Meadow wood", VI: "Bình Địa Mộc"},
	"Adobe earth":      {Code: "nayin.adobe_earth", EN: "Adobe earth", VI: "Bích Thượng Thổ"},
	"Precious metal":   {Code: "nayin.precious_metal", EN: "Precious metal", VI: "Kim Bạch Kim"},
	"Lamp fire":        {Code: "nayin.lamp_fire", EN: "Lamp fire", VI: "Phúc Đăng Hỏa"},
	"Sky water":        {Code: "nayin.sky_water", EN: "Sky water", VI: "Thiên Hà Thủy"},
	"Highway earth":    {Code: "nayin.highway_earth", EN: "Highway earth", VI: "Đại Trạch Thổ"},
	"Jewellery metal":  {Code: "nayin.jewellery_metal", EN: "Jewellery metal", VI: "Thoa Xuyến Kim"},
	"Mulberry wood":    {Code: "nayin.mulberry_wood", EN: "Mulberry wood", VI: "Tang Đố Mộc"},
	"Rapids water":     {Code: "nayin.rapids_water", EN: "Rapids water", VI: "Đại Khê Thủy"},
	"Desert earth":     {Code: "nayin.desert_earth", EN: "Desert earth", VI: "Sa Trung Thổ"},
	"Sun fire":         {Code: "nayin.sun_fire", EN: "Sun fire", VI: "Thiên Thượng Hỏa"},
	"Pomegranate wood": {Code: "nayin.pomegranate_wood", EN: "Pomegranate wood", VI: "Thạch Lựu Mộc"},
	"Ocean water":      {Code: "nayin.ocean_water", EN: "Ocean water", VI: "Đại Hải Thủy"},
}
