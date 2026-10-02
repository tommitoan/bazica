package analysis

// Branches derived from the Loc branch of the Day Master stem.
var (
	amLocTable = deriveStemTable(locTable, func(_ int, b []int) []int { return []int{harmony(b[0])} })
	// giapLocTable lists the two branches either side of Loc.
	giapLocTable   = deriveStemTable(locTable, func(_ int, b []int) []int { return []int{mod12(b[0] - 1), mod12(b[0] + 1)} })
	duongNhanTable = deriveStemTable(locTable, func(_ int, b []int) []int { return []int{mod12(b[0] + 1)} })
	phiNhanTable   = deriveStemTable(locTable, func(_ int, b []int) []int { return []int{opposite(mod12(b[0] + 1))} })
	// kinhDuongTable is one step after Loc for Yang stems and one step before for Yin stems.
	kinhDuongTable = deriveStemTable(locTable, func(stem int, b []int) []int {
		if stem%2 == 0 {
			return []int{mod12(b[0] + 1)}
		}
		return []int{mod12(b[0] - 1)}
	})

	hongLoanBranches = deriveBranchTable(func(b int) []int { return []int{hongLoanTable[b]} })
	thienHyBranches  = deriveBranchTable(func(b int) []int { return []int{opposite(hongLoanTable[b])} })
	// thienYBranches is the branch before the month branch.
	thienYBranches = deriveBranchTable(func(b int) []int { return []int{mod12(b - 1)} })
)

// starRules is the star table in registry order. Each row is one star; the
// evaluator in stars.go applies every row to every pillar it may land in.
var starRules = []starRule{
	byStem("thien_at_quy_nhan", dayOrYear, thienAtTable, inAll),
	byStem("quoc_an", dayOrYear, quocAnTable, inAll),
	byStem("van_xuong", dayOrYear, vanXuongTable, inAll),
	byStem("phuc_tinh", dayOrYear, phucTinhTable, inAll),
	byStem("thai_cuc_quy_nhan", yearOnly, thaiCucTable, inAll),
	byStem("tue_loc", dayOnly, locTable, onlyPillar(yearPillar)),
	byStem("kien_loc", dayOnly, locTable, onlyPillar(monthPillar)),
	byStem("toa_loc", dayOnly, locTable, onlyPillar(dayPillar)),
	byStem("quy_loc", dayOnly, locTable, onlyPillar(hourPillar)),
	byStem("am_loc", dayOnly, amLocTable, inAll),
	byStem("giap_loc", dayOnly, giapLocTable, inAll),
	byStem("duong_nhan", dayOnly, duongNhanTable, inAll),
	byStem("phi_nhan", dayOnly, phiNhanTable, inAll),
	byStem("kinh_duong", dayOnly, kinhDuongTable, inAll),
	byStem("kim_du", dayOnly, kimDuTable, inAll),
	byStem("hoc_duong", dayOnly, hocDuongTable, inMonthHour),
	byGroup("vong_than", yearOrDay, vongThanTargets, inAll),
	byGroup("kiep_sat", yearOrDay, kiepSatTargets, inAll),
	byGroup("dich_ma", yearOrDay, dichMaTargets, inAll),
	byGroup("tuong_tinh", yearOrDay, tuongTinhTargets, inAll),
	byGroup("hoa_cai", yearOrDay, hoaCaiTargets, inAll),
	byGroup("dao_hoa", yearOrDay, daoHoaTargets, inAll),
	byGroup("tai_sat", yearOnly, taiSatTargets, inAll),
	byBranch("hong_loan", yearPillar, hongLoanBranches, inAll),
	byBranch("thien_hy", yearPillar, thienHyBranches, inAll),
	giaoCau("giao", false),
	giaoCau("cau", true),
	nguyetDuc("nguyet_duc_quy_nhan"),
	thienDuc("thien_duc_quy_nhan"),
	byBranch("thien_y", monthPillar, thienYBranches, inAll),
	byPair("thap_ac_dai_bai", thapAcPairs, onlyPillar(dayPillar)),
	byPair("khoi_cuong", khoiCuongPairs, inAll),
	byPair("tien_than", tienThanPairs, inAll),
	byPair("kim_than", kimThanPairs, inDayHour),
	byPair("am_duong_sai_thac", amDuongSaiThacPairs, onlyPillar(dayPillar)),
	byPair("am_duong_sat", amDuongSatPairs, onlyPillar(dayPillar)),
	byPair("nhat_quy", nhatQuyPairs, onlyPillar(dayPillar)),
	byPair("luc_tu", lucTuPairs, onlyPillar(dayPillar)),
	voidStar("khong_vong"),
}
