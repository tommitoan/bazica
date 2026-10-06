package analysis

import (
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/model"
)

// pillarsHaving returns the pillar positions of an evaluation that carry a star.
func pillarsHaving(stars [4][]model.LocalizedTerm, code string) []int {
	var out []int
	for p := range stars {
		if has(codesOf(stars, p), code) {
			out = append(out, p)
		}
	}
	return out
}

func sameInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestDaiHaoFollowsTheDirectionOfGiaoAndCau(t *testing.T) {
	// Yi-Hai year: the year branch Hai (11) clashes with Si (5). A Yin-year woman moves one step on (Wu, 6);
	// a Yin-year man moves one step back (Chen, 4).
	for _, tc := range []struct {
		male bool
		want int
	}{{false, 6}, {true, 4}} {
		stars := evaluateStars(facts([4]int{1, 8, 6, 8}, [4]int{11, tc.want, 0, 2}, tc.male))
		if got := pillarsHaving(stars, "dai_hao"); !sameInts(got, []int{monthPillar}) {
			t.Errorf("male=%v: dai_hao on %v, want only the month pillar", tc.male, got)
		}
	}
	// The other direction must not award it: a woman does not get Chen from a Yin year.
	stars := evaluateStars(facts([4]int{1, 8, 6, 8}, [4]int{11, 4, 0, 2}, false))
	if got := pillarsHaving(stars, "dai_hao"); len(got) != 0 {
		t.Errorf("Yin-year woman with Chen: dai_hao on %v, want none", got)
	}
	// A Yang-year man (Jia-Wu): clash of Wu is Zi (0), one step on is Chou (1); a Yang-year woman gets Hai (11).
	man := evaluateStars(facts([4]int{0, 8, 6, 8}, [4]int{6, 1, 5, 2}, true))
	woman := evaluateStars(facts([4]int{0, 8, 6, 8}, [4]int{6, 11, 5, 2}, false))
	if !sameInts(pillarsHaving(man, "dai_hao"), []int{monthPillar}) || !sameInts(pillarsHaving(woman, "dai_hao"), []int{monthPillar}) {
		t.Error("Yang year: man gets the branch after the clash, woman the branch before it")
	}
}

func TestThienLaAndDiaVongNeedTheDayBranchAsPartner(t *testing.T) {
	// Day branch Chen (4): Si (5) on the month pillar earns Thien La; the day pillar itself does not.
	stars := evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{4, 5, 4, 2}, true))
	if got := pillarsHaving(stars, "thien_la"); !sameInts(got, []int{monthPillar}) {
		t.Errorf("thien_la on %v, want the month pillar only (a Chen year pillar has the same branch as the day)", got)
	}
	// The partner must be the day branch: Si in the year pillar and Chen in the month pillar earn nothing.
	stars = evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{5, 4, 8, 2}, true))
	if got := pillarsHaving(stars, "thien_la"); len(got) != 0 {
		t.Errorf("thien_la on %v without a Chen/Si day branch, want none", got)
	}
	// Dia Vong: day branch Hai (11), Xu (10) on the hour pillar.
	stars = evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{2, 3, 11, 10}, true))
	if got := pillarsHaving(stars, "dia_vong"); !sameInts(got, []int{hourPillar}) {
		t.Errorf("dia_vong on %v, want the hour pillar", got)
	}
}

func TestTamKyNeedsTheStemsInOrderOnConsecutivePillars(t *testing.T) {
	// Ren, Gui, Xin on the month, day and hour pillars.
	forward := evaluateStars(facts([4]int{0, 8, 9, 7}, [4]int{0, 0, 0, 0}, true))
	if got := pillarsHaving(forward, "tam_ky_trung"); !sameInts(got, []int{monthPillar, dayPillar, hourPillar}) {
		t.Errorf("tam_ky_trung on %v, want month, day and hour", got)
	}
	// The reverse order does not count.
	reversed := evaluateStars(facts([4]int{0, 7, 9, 8}, [4]int{0, 0, 0, 0}, true))
	if got := pillarsHaving(reversed, "tam_ky_trung"); len(got) != 0 {
		t.Errorf("reversed stems: tam_ky_trung on %v, want none", got)
	}
	// Pillars that are not consecutive do not count.
	gap := evaluateStars(facts([4]int{8, 9, 0, 7}, [4]int{0, 0, 0, 0}, true))
	if got := pillarsHaving(gap, "tam_ky_trung"); len(got) != 0 {
		t.Errorf("stems with a gap: tam_ky_trung on %v, want none", got)
	}
	// Jia Wu Geng is Upper and Yi Bing Ding is Lower; the page labels them that way.
	upper := evaluateStars(facts([4]int{0, 4, 6, 3}, [4]int{0, 0, 0, 0}, true))
	lower := evaluateStars(facts([4]int{1, 2, 3, 0}, [4]int{0, 0, 0, 0}, true))
	if !sameInts(pillarsHaving(upper, "tam_ky_thuong"), []int{yearPillar, monthPillar, dayPillar}) ||
		!sameInts(pillarsHaving(lower, "tam_ky_ha"), []int{yearPillar, monthPillar, dayPillar}) {
		t.Error("Upper is Jia Wu Geng and Lower is Yi Bing Ding, on the year, month and day pillars")
	}
}

func TestCoThanAndQuaTuFollowTheSeasonOfTheYearBranch(t *testing.T) {
	// Tiger year (Dan-Mao-Chen group): Lonely Spirit at Ox (1), Widowhood at Snake (5); never on the year pillar.
	stars := evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{2, 1, 5, 1}, true))
	if got := pillarsHaving(stars, "co_than"); !sameInts(got, []int{monthPillar, hourPillar}) {
		t.Errorf("co_than on %v, want month and hour", got)
	}
	if got := pillarsHaving(stars, "qua_tu"); !sameInts(got, []int{dayPillar}) {
		t.Errorf("qua_tu on %v, want day", got)
	}
	// Pig year (Hai-Zi-Chou group): Lonely Spirit at Tiger, Widowhood at Dog.
	pig := evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{11, 2, 10, 0}, true))
	if !sameInts(pillarsHaving(pig, "co_than"), []int{monthPillar}) || !sameInts(pillarsHaving(pig, "qua_tu"), []int{dayPillar}) {
		t.Error("Pig year: Lonely Spirit at Tiger and Widowhood at Dog")
	}
}

func TestHongDiemUsesTheDayOrTheYearStemAndHocSyOnlyTheDayStem(t *testing.T) {
	// Day stem Geng (6): Red Allure at Xu (10). Year stem Jia (0): Wu (6) and Shen (8).
	byDay := evaluateStars(facts([4]int{3, 3, 6, 3}, [4]int{2, 10, 4, 4}, true))
	if got := pillarsHaving(byDay, "hong_diem"); !sameInts(got, []int{monthPillar}) {
		t.Errorf("day stem Geng: hong_diem on %v, want the month pillar", got)
	}
	byYear := evaluateStars(facts([4]int{0, 3, 3, 3}, [4]int{2, 8, 4, 4}, true))
	if got := pillarsHaving(byYear, "hong_diem"); !sameInts(got, []int{monthPillar}) {
		t.Errorf("year stem Jia: hong_diem on %v, want the month pillar (Shen)", got)
	}
	// Scholar follows the day stem only: day stem Geng has Wu (6); a Geng year stem with a Jia day stem does not.
	stars := evaluateStars(facts([4]int{6, 3, 0, 3}, [4]int{2, 6, 4, 4}, true))
	if got := pillarsHaving(stars, "hoc_sy"); len(got) != 0 {
		t.Errorf("year stem Geng must not award hoc_sy, got %v", got)
	}
	stars = evaluateStars(facts([4]int{0, 3, 6, 3}, [4]int{2, 6, 4, 4}, true))
	if got := pillarsHaving(stars, "hoc_sy"); !sameInts(got, []int{monthPillar}) {
		t.Errorf("day stem Geng: hoc_sy on %v, want the month pillar", got)
	}
}

func TestDucAndTuQuyNhanFollowTheMonthBranchGroup(t *testing.T) {
	// Month Zi (Shen-Zi-Chen group): Virtue Noble on stems Wu, Ji, Ren, Gui; no Elegance Noble.
	stars := evaluateStars(facts([4]int{4, 0, 9, 2}, [4]int{0, 0, 0, 0}, true))
	if got := pillarsHaving(stars, "duc_quy_nhan"); !sameInts(got, []int{yearPillar, dayPillar}) {
		t.Errorf("duc_quy_nhan on %v, want year and day", got)
	}
	if got := pillarsHaving(stars, "tu_quy_nhan"); len(got) != 0 {
		t.Errorf("tu_quy_nhan on %v, want none in the Shen-Zi-Chen group", got)
	}
	// Month Yin (Yin-Wu-Xu group): Elegance Noble on Bing and Ding, no Virtue Noble.
	stars = evaluateStars(facts([4]int{2, 0, 3, 0}, [4]int{0, 2, 0, 0}, true))
	if got := pillarsHaving(stars, "tu_quy_nhan"); !sameInts(got, []int{yearPillar, dayPillar}) {
		t.Errorf("tu_quy_nhan on %v, want year and day", got)
	}
	if got := pillarsHaving(stars, "duc_quy_nhan"); len(got) != 0 {
		t.Errorf("duc_quy_nhan on %v, want none in the Yin-Wu-Xu group", got)
	}
}

func TestSeasonalDayPillarStars(t *testing.T) {
	// Spring (month Mao, 3) with the Yi-Mao day: Heaven Earth Turning Sha and Four Wasted (Yi-Mao is autumn's) split by season.
	spring := evaluateStars(facts([4]int{0, 0, 1, 0}, [4]int{0, 3, 3, 0}, true))
	if !has(codesOf(spring, dayPillar), "thien_dia_chuyen_sat") || has(codesOf(spring, dayPillar), "tu_phe") {
		t.Error("spring Yi-Mao day: Thien Dia Chuyen Sat only")
	}
	autumn := evaluateStars(facts([4]int{0, 0, 1, 0}, [4]int{0, 9, 3, 0}, true))
	if !has(codesOf(autumn, dayPillar), "tu_phe") || has(codesOf(autumn, dayPillar), "thien_dia_chuyen_sat") {
		t.Error("autumn Yi-Mao day: Tu Phe only")
	}
	// Heavenly Pardon: winter (month Zi) with a Jia-Zi day; the same day in summer earns nothing.
	winter := evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{0, 0, 0, 0}, true))
	summer := evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{0, 6, 0, 0}, true))
	if !has(codesOf(winter, dayPillar), "thien_xich_quy") || has(codesOf(summer, dayPillar), "thien_xich_quy") {
		t.Error("Jia-Zi day carries Thien Xich Quy in winter only")
	}
	// The stars belong to the day pillar only.
	if got := pillarsHaving(winter, "thien_xich_quy"); !sameInts(got, []int{dayPillar}) {
		t.Errorf("thien_xich_quy on %v, want only the day pillar", got)
	}
}

func TestCuuQuyPhongHaiMarksTheHourPillarToo(t *testing.T) {
	// Ren-Wu (8, 6) is on the list; as the hour pillar it carries the star, as the month pillar it does not.
	stars := evaluateStars(facts([4]int{0, 8, 0, 8}, [4]int{0, 6, 0, 6}, true))
	if got := pillarsHaving(stars, "cuu_quy_phong_hai"); !sameInts(got, []int{hourPillar}) {
		t.Errorf("cuu_quy_phong_hai on %v, want only the hour pillar here", got)
	}
}

func TestMoonGeneralNeedsABirthDate(t *testing.T) {
	// Born 2022-06-21 (the Summer Solstice day): the general is Wei (7). Put Wei on the year pillar.
	pairs := [4][2]int{{2, 6}, {2, 2}, {2, 8}, {2, 2}} // public values: Yi-Wei, Yi-Mao, Yi-You, Yi-Mao
	chart := chartOf(pairs)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	if has(starCodes(chart.FourPillar.YearPillar.Analysis.Stars), "thai_duong") {
		t.Error("Attach does not know the birth date and must not award the Moon General")
	}
	chart = chartOf(pairs)
	birth := time.Date(2022, 6, 21, 15, 30, 0, 0, time.FixedZone("ICT", 7*3600))
	if err := AttachAt(chart, model.GenderMale, birth); err != nil {
		t.Fatal(err)
	}
	if !has(starCodes(chart.FourPillar.YearPillar.Analysis.Stars), "thai_duong") {
		t.Error("AttachAt should award the Moon General to the pillar holding Wei")
	}
}
