package analysis

import (
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/model"
)

func starCodes(terms []model.LocalizedTerm) []string {
	codes := make([]string, 0, len(terms))
	for _, t := range terms {
		codes = append(codes, strings.TrimPrefix(t.Code, "star."))
	}
	return codes
}

// codesOf returns the star codes of one pillar of an evaluation, never nil.
func codesOf(stars [4][]model.LocalizedTerm, pillar int) []string { return starCodes(stars[pillar]) }

func facts(stems, branches [4]int, male bool) *chartFacts {
	return &chartFacts{stem: stems, branch: branches, male: male}
}

func has(codes []string, code string) bool {
	for _, c := range codes {
		if c == code {
			return true
		}
	}
	return false
}

type starCase struct {
	Name    string    `json:"name"`
	Pillars [4][2]int `json:"pillars"`
	Male    bool      `json:"male"`
	// Birth is the wall-clock birth time in Vietnam; when set, the Moon General is awarded.
	Birth string `json:"birth"`
	// Skip lists star codes the comparison ignores for this chart (see the generator of the fixture).
	Skip  []string    `json:"skip"`
	Stars [4][]string `json:"stars"`
}

// without drops the skipped star codes from every pillar of a result.
func without(stars [4][]string, skip []string) [4][]string {
	var out [4][]string
	for i, pillar := range stars {
		out[i] = []string{}
		for _, code := range pillar {
			if !slices.Contains(skip, code) {
				out[i] = append(out[i], code)
			}
		}
	}
	return out
}

func loadStarCases(t *testing.T, file string) []starCase {
	t.Helper()
	data, err := os.ReadFile("../../testdata/analysis/" + file)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var f struct {
		Charts []starCase `json:"charts"`
	}
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	return f.Charts
}

func attachedStars(t *testing.T, c starCase) [4][]string {
	t.Helper()
	gender := model.GenderFemale
	if c.Male {
		gender = model.GenderMale
	}
	chart := chartOf(c.Pillars)
	var err error
	if c.Birth != "" {
		var birth time.Time
		birth, err = time.ParseInLocation("2006-01-02T15:04", c.Birth, time.FixedZone("ICT", 7*3600))
		if err != nil {
			t.Fatalf("birth %q: %v", c.Birth, err)
		}
		err = AttachAt(chart, gender, birth)
	} else {
		err = Attach(chart, gender)
	}
	if err != nil {
		t.Fatalf("Attach(%v): %v", c.Pillars, err)
	}
	fp := chart.FourPillar
	return [4][]string{
		starCodes(fp.YearPillar.Analysis.Stars), starCodes(fp.MonthPillar.Analysis.Stars),
		starCodes(fp.DayPillar.Analysis.Stars), starCodes(fp.HourPillar.Analysis.Stars),
	}
}

func TestStarRulesFollowTheSharedRegistry(t *testing.T) {
	data, err := os.ReadFile("../../testdata/analysis/term_registry.json")
	if err != nil {
		t.Skipf("registry fixture not found: %v", err)
	}
	var file registryFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	var want []model.LocalizedTerm
	for _, term := range file.Terms {
		if strings.HasPrefix(term.Code, "star.") {
			want = append(want, term)
		}
	}
	if len(starRules) != len(want) || len(starTerms) != len(want) {
		t.Fatalf("%d rules and %d terms for %d registry stars", len(starRules), len(starTerms), len(want))
	}
	for i, term := range want {
		code := strings.TrimPrefix(term.Code, "star.")
		if starRules[i].code != code {
			t.Errorf("rule %d is %q, registry order has %q", i, starRules[i].code, code)
		}
		if starTerms[code] != term {
			t.Errorf("%s: term %+v, want %+v", code, starTerms[code], term)
		}
	}
}

// Charts recorded from the reference page: every star it printed must be reproduced exactly.
func TestStarsMatchTheChartsThePagePrinted(t *testing.T) {
	cases := loadStarCases(t, "stars_observed.json")
	if len(cases) < 250 {
		t.Fatalf("only %d observed charts", len(cases))
	}
	for _, c := range cases {
		got, want := without(attachedStars(t, c), c.Skip), without(c.Stars, c.Skip)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %v, want %v", c.Name, got, want)
		}
	}
}

// Random valid charts, with the stars an independent oracle awards. They
// carry no birth date, so the Moon General is not part of them.
func TestStarsMatchTheOracleOnRandomCharts(t *testing.T) {
	for i, c := range loadStarCases(t, "stars_oracle_cases.json") {
		if got := attachedStars(t, c); !reflect.DeepEqual(got, c.Stars) {
			t.Fatalf("case %d %v male=%v: got %v, want %v", i, c.Pillars, c.Male, got, c.Stars)
		}
	}
}

func TestEveryRegistryStarCanBeAwarded(t *testing.T) {
	seen := map[string]bool{}
	for _, file := range []string{"stars_oracle_cases.json", "stars_observed.json"} {
		for _, c := range loadStarCases(t, file) {
			for _, pillar := range attachedStars(t, c) {
				for _, code := range pillar {
					seen[code] = true
				}
			}
		}
	}
	for _, rule := range starRules {
		if !seen[rule.code] {
			t.Errorf("star %s was never awarded in the random cases", rule.code)
		}
	}
}

func TestStarsAreNeverNilAndSerialiseAsAnEmptyList(t *testing.T) {
	// The oracle awards no star to the year pillar of this chart (public stem and branch values).
	chart := chartOf([4][2]int{{7, 7}, {9, 11}, {9, 7}, {3, 3}})
	if err := Attach(chart, model.GenderFemale); err != nil {
		t.Fatal(err)
	}
	year := chart.FourPillar.YearPillar.Analysis
	if year.Stars == nil || len(year.Stars) != 0 {
		t.Fatalf("year stars = %#v, want an empty non-nil list", year.Stars)
	}
	data, err := json.Marshal(year)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"stars":[]`) {
		t.Errorf("JSON = %s, want stars as []", data)
	}
}

func TestOnlyRegistryStarsAreEverEmitted(t *testing.T) {
	for _, c := range loadStarCases(t, "stars_oracle_cases.json") {
		for _, pillar := range attachedStars(t, c) {
			for _, code := range pillar {
				if _, ok := starTerms[code]; !ok {
					t.Fatalf("emitted unregistered star %q", code)
				}
			}
		}
	}
}

func TestLocIsNamedByThePillarItLandsIn(t *testing.T) {
	// Day stem Jia has Loc at Tiger (index 2): put Tiger in each pillar in turn.
	for pillar, want := range []string{"tue_loc", "kien_loc", "toa_loc", "quy_loc"} {
		branches := [4]int{0, 4, 6, 8}
		branches[pillar] = 2
		stars := evaluateStars(facts([4]int{0, 0, 0, 0}, branches, true))
		for p := 0; p < 4; p++ {
			if got := has(codesOf(stars, p), want); got != (p == pillar) {
				t.Errorf("%s in pillar %d = %v, want %v", want, p, got, p == pillar)
			}
		}
	}
}

func TestGiaoAndCauFollowYearPolarityAndGender(t *testing.T) {
	// Rat year: three ahead is Rabbit (3), three behind is Rooster (9).
	atRabbit := facts([4]int{0, 0, 0, 0}, [4]int{0, 3, 4, 6}, true)
	atRooster := facts([4]int{0, 0, 0, 0}, [4]int{0, 9, 4, 6}, true)
	if !has(codesOf(evaluateStars(atRabbit), monthPillar), "cau") || has(codesOf(evaluateStars(atRabbit), monthPillar), "giao") {
		t.Error("Yang-year man: Cau should sit three ahead")
	}
	if !has(codesOf(evaluateStars(atRooster), monthPillar), "giao") {
		t.Error("Yang-year man: Giao should sit three behind")
	}
	atRabbit.male, atRooster.male = false, false
	if !has(codesOf(evaluateStars(atRabbit), monthPillar), "giao") || !has(codesOf(evaluateStars(atRooster), monthPillar), "cau") {
		t.Error("Yang-year woman: Giao and Cau swap")
	}
	// A Yin year (Yi, index 1) and a woman behaves like a Yang year and a man.
	yinWoman := facts([4]int{1, 0, 0, 0}, [4]int{0, 3, 4, 6}, false)
	if !has(codesOf(evaluateStars(yinWoman), monthPillar), "cau") {
		t.Error("Yin-year woman: Cau should sit three ahead")
	}
}

func TestHocDuongOnlyOnMonthAndHourPillars(t *testing.T) {
	// Day stem Jia: Study Hall at Pig (11).
	stars := evaluateStars(facts([4]int{1, 1, 0, 1}, [4]int{11, 11, 11, 11}, true))
	for p, want := range [4]bool{false, true, false, true} {
		if got := has(codesOf(stars, p), "hoc_duong"); got != want {
			t.Errorf("hoc_duong in pillar %d = %v, want %v", p, got, want)
		}
	}
}

func TestKimThanOnlyOnDayAndHourPillars(t *testing.T) {
	// Yi-Chou (1, 1) carries Golden Spirit, but only when it is the day or hour pillar.
	stars := evaluateStars(facts([4]int{1, 1, 1, 1}, [4]int{1, 1, 1, 1}, true))
	for p, want := range [4]bool{false, false, true, true} {
		if got := has(codesOf(stars, p), "kim_than"); got != want {
			t.Errorf("kim_than in pillar %d = %v, want %v", p, got, want)
		}
	}
}

func TestVoidNeverMarksTheDayPillar(t *testing.T) {
	// Jia-Zi day: Xu (10) and Hai (11) are void.
	stars := evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{10, 11, 0, 10}, true))
	for p, want := range [4]bool{true, true, false, true} {
		if got := has(codesOf(stars, p), "khong_vong"); got != want {
			t.Errorf("khong_vong in pillar %d = %v, want %v", p, got, want)
		}
	}
}

func TestThaiCucUsesTheYearStemOnly(t *testing.T) {
	// Year stem Jia: Supreme Noble at Rat (0) and Horse (6). With a Jia day stem but a Yi... use Bing year stem.
	withJiaYear := evaluateStars(facts([4]int{0, 2, 2, 2}, [4]int{0, 4, 4, 4}, true))
	if !has(codesOf(withJiaYear, yearPillar), "thai_cuc_quy_nhan") {
		t.Error("Jia year stem with a Rat pillar should earn Thai Cuc")
	}
	withBingYear := evaluateStars(facts([4]int{2, 0, 0, 0}, [4]int{4, 0, 6, 4}, true))
	if has(codesOf(withBingYear, monthPillar), "thai_cuc_quy_nhan") {
		t.Error("Thai Cuc follows the year stem; a Jia day stem must not award it")
	}
}

func TestPillarIdentityStarsNeedTheExactPair(t *testing.T) {
	// Bing-Zi (2, 0) as the day pillar carries Yin Yang Error and Yin Yang Sha; Bing-Chou only neither of those.
	both := evaluateStars(facts([4]int{0, 0, 2, 0}, [4]int{0, 0, 0, 0}, true))
	if !has(codesOf(both, dayPillar), "am_duong_sai_thac") || !has(codesOf(both, dayPillar), "am_duong_sat") {
		t.Error("Bing-Zi day pillar should carry both Yin Yang stars")
	}
	none := evaluateStars(facts([4]int{0, 0, 3, 0}, [4]int{0, 0, 1, 0}, true))
	if has(codesOf(none, dayPillar), "am_duong_sat") {
		t.Error("Ding-Chou must not carry Yin Yang Sha")
	}
	notDay := evaluateStars(facts([4]int{2, 0, 0, 0}, [4]int{0, 0, 4, 0}, true))
	if has(codesOf(notDay, yearPillar), "am_duong_sat") {
		t.Error("Yin Yang Sha belongs to the day pillar only")
	}
}

func TestDerivedLocStarsFollowTheLocBranch(t *testing.T) {
	// Day stem Jia: Loc at Tiger (2). Hidden Loc is its harmony partner Pig (11), Duong Nhan sits
	// at Rabbit (3), Phi Nhan opposite it at Rooster (9), and Giap Loc at Ox (1) and Rabbit (3).
	stars := evaluateStars(facts([4]int{0, 0, 0, 0}, [4]int{11, 3, 9, 1}, true))
	for pillar, want := range map[int][]string{
		yearPillar:  {"am_loc"},
		monthPillar: {"giap_loc", "duong_nhan", "kinh_duong"},
		dayPillar:   {"phi_nhan"},
		hourPillar:  {"giap_loc"},
	} {
		for _, code := range want {
			if !has(codesOf(stars, pillar), code) {
				t.Errorf("pillar %d is missing %s: %v", pillar, code, codesOf(stars, pillar))
			}
		}
	}
	// A Yin day stem puts Kinh Duong one step before Loc: Yi has Loc at Rabbit (3), Kinh Duong at Tiger (2).
	yin := evaluateStars(facts([4]int{0, 0, 1, 0}, [4]int{0, 2, 0, 0}, true))
	if !has(codesOf(yin, monthPillar), "kinh_duong") {
		t.Error("Yin day stem: Kinh Duong should sit one step before Loc")
	}
}
