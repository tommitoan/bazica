package analysis

import "github.com/tommitoan/bazica/model"

// starNature classifies each star, in registry order, as auspicious (cát), inauspicious (hung)
// or mixed (tùy cục: the effect depends on the rest of the chart). Schools differ, so the value
// is a hint for display. Every star in starTerms must appear here (see star_nature_test.go).
var starNature = map[string]string{
	"thien_at_quy_nhan":    model.NatureAuspicious,
	"quoc_an":              model.NatureAuspicious,
	"van_xuong":            model.NatureAuspicious,
	"phuc_tinh":            model.NatureAuspicious,
	"thai_cuc_quy_nhan":    model.NatureAuspicious,
	"tue_loc":              model.NatureAuspicious,
	"kien_loc":             model.NatureAuspicious,
	"toa_loc":              model.NatureAuspicious,
	"quy_loc":              model.NatureAuspicious,
	"am_loc":               model.NatureMixed,
	"giap_loc":             model.NatureAuspicious,
	"duong_nhan":           model.NatureMixed,
	"phi_nhan":             model.NatureInauspicious,
	"kinh_duong":           model.NatureMixed,
	"kim_du":               model.NatureAuspicious,
	"hoc_duong":            model.NatureAuspicious,
	"vong_than":            model.NatureInauspicious,
	"kiep_sat":             model.NatureInauspicious,
	"dich_ma":              model.NatureMixed,
	"tuong_tinh":           model.NatureAuspicious,
	"hoa_cai":              model.NatureMixed,
	"dao_hoa":              model.NatureMixed,
	"tai_sat":              model.NatureInauspicious,
	"hong_loan":            model.NatureAuspicious,
	"thien_hy":             model.NatureAuspicious,
	"giao":                 model.NatureInauspicious,
	"cau":                  model.NatureInauspicious,
	"nguyet_duc_quy_nhan":  model.NatureAuspicious,
	"thien_duc_quy_nhan":   model.NatureAuspicious,
	"thien_y":              model.NatureAuspicious,
	"thap_ac_dai_bai":      model.NatureInauspicious,
	"khoi_cuong":           model.NatureMixed,
	"tien_than":            model.NatureAuspicious,
	"kim_than":             model.NatureMixed,
	"am_duong_sai_thac":    model.NatureInauspicious,
	"am_duong_sat":         model.NatureInauspicious,
	"nhat_quy":             model.NatureAuspicious,
	"luc_tu":               model.NatureAuspicious,
	"khong_vong":           model.NatureMixed,
	"hong_diem":            model.NatureMixed,
	"hoc_sy":               model.NatureAuspicious,
	"duc_quy_nhan":         model.NatureAuspicious,
	"tu_quy_nhan":          model.NatureAuspicious,
	"thai_duong":           model.NatureAuspicious,
	"tam_ky_thuong":        model.NatureAuspicious,
	"tam_ky_trung":         model.NatureAuspicious,
	"tam_ky_ha":            model.NatureAuspicious,
	"dai_hao":              model.NatureInauspicious,
	"co_than":              model.NatureInauspicious,
	"qua_tu":               model.NatureInauspicious,
	"thien_la":             model.NatureInauspicious,
	"dia_vong":             model.NatureInauspicious,
	"co_loan_sat":          model.NatureInauspicious,
	"nhat_nhan":            model.NatureInauspicious,
	"cuu_quy_phong_hai":    model.NatureInauspicious,
	"thien_dia_chuyen_sat": model.NatureInauspicious,
	"tu_phe":               model.NatureInauspicious,
	"thien_xich_quy":       model.NatureAuspicious,
}

// init puts each star's nature into its registry term, so every place that emits a star carries it.
func init() {
	for code, nature := range starNature {
		term := starTerms[code]
		term.Nature = nature
		starTerms[code] = term
	}
}
