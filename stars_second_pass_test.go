package bazica

import (
	"strings"
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

func starCodesOf(analysis *model.PillarAnalysis) []string {
	var codes []string
	for _, s := range analysis.Stars {
		codes = append(codes, strings.TrimPrefix(s.Code, "star."))
	}
	return codes
}

func containsCode(codes []string, code string) bool {
	for _, c := range codes {
		if c == code {
			return true
		}
	}
	return false
}

// A woman born 1995-06-08 11:01 (Yi-Hai, Ren-Wu, Geng-Wu, Ren-Wu) has Great Depletion on each Wu pillar;
// the man born at the same moment does not, because the direction follows gender.
func TestDaiHaoDependsOnGender(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	birth := time.Date(1995, 6, 8, 11, 1, 0, 0, loc)

	woman, err := GetBaziChart(birth, loc, model.GenderFemale)
	if err != nil {
		t.Fatal(err)
	}
	fp := woman.FourPillar
	for name, a := range map[string]*model.PillarAnalysis{"month": fp.MonthPillar.Analysis, "day": fp.DayPillar.Analysis, "hour": fp.HourPillar.Analysis} {
		if !containsCode(starCodesOf(a), "dai_hao") {
			t.Errorf("woman, %s pillar: no dai_hao in %v", name, starCodesOf(a))
		}
	}
	if containsCode(starCodesOf(fp.YearPillar.Analysis), "dai_hao") {
		t.Error("woman, year pillar (Hai): must not have dai_hao")
	}

	man, err := GetBaziChart(birth, loc, model.GenderMale)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range []*model.PillarAnalysis{man.FourPillar.YearPillar.Analysis, man.FourPillar.MonthPillar.Analysis, man.FourPillar.DayPillar.Analysis, man.FourPillar.HourPillar.Analysis} {
		if containsCode(starCodesOf(a), "dai_hao") {
			t.Errorf("man: unexpected dai_hao in %v", starCodesOf(a))
		}
	}
}

// The Moon General is awarded to exactly the pillars whose branch is the general of the birth date.
func TestMoonGeneralFollowsTheBirthDate(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	births := []time.Time{
		time.Date(1995, 6, 8, 11, 1, 0, 0, loc),
		time.Date(2022, 6, 21, 15, 30, 0, 0, loc),
		time.Date(2023, 12, 22, 3, 30, 0, 0, loc),
		time.Date(1984, 3, 3, 1, 50, 0, 0, loc),
		time.Date(2010, 9, 14, 22, 0, 0, 0, loc),
	}
	for _, birth := range births {
		chart, err := GetBaziChart(birth, loc, model.GenderMale)
		if err != nil {
			t.Fatal(err)
		}
		general, err := utils.GetMoonGeneral(birth)
		if err != nil {
			t.Fatal(err)
		}
		fp := chart.FourPillar
		pillars := []struct {
			branch   int
			analysis *model.PillarAnalysis
		}{
			{fp.YearPillar.EarthlyBranch.Value, fp.YearPillar.Analysis}, {fp.MonthPillar.EarthlyBranch.Value, fp.MonthPillar.Analysis},
			{fp.DayPillar.EarthlyBranch.Value, fp.DayPillar.Analysis}, {fp.HourPillar.EarthlyBranch.Value, fp.HourPillar.Analysis},
		}
		for i, p := range pillars {
			// Public branch values start at Tiger = 1; the general is a Rat-first index.
			want := (p.branch+1)%12 == general
			if got := containsCode(starCodesOf(p.analysis), "thai_duong"); got != want {
				t.Errorf("%v pillar %d: thai_duong = %v, want %v (general %d, branch value %d)", birth, i, got, want, general, p.branch)
			}
		}
	}
}
