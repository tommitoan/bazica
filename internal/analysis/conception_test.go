package analysis

import (
	"testing"

	"github.com/tommitoan/bazica/v2/internal/utils"
)

func TestThaiNguyenAndThaiTucKnownChart(t *testing.T) {
	// 1995-06-08 reference chart: month Ren-Wu, day Geng-Wu.
	if s, b := thaiNguyen(9, 5); s != 10 || b != 8 {
		t.Errorf("thaiNguyen(Ren, Wu) = %d/%d, want Gui/You (10/8)", s, b)
	}
	if s, b := thaiTuc(7, 5); s != 2 || b != 6 {
		t.Errorf("thaiTuc(Geng, Wu) = %d/%d, want Yi/Wei (2/6)", s, b)
	}
}

func TestLifePalaceKnownCharts(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		yearStem, month, hour int
		wantStem, wantBranch  int
	}{
		{"1995 chart", 2, 5, 5, 8, 4},
		{"Rat palace, Xin year", 8, 1, 2, 5, 11},
		{"Ox palace, Xin year", 8, 1, 1, 6, 12},
		{"Rat palace, Ren year", 9, 1, 2, 7, 11},
		{"Ox palace, Ren year", 9, 1, 1, 8, 12},
	} {
		s, b := lifePalace(tc.yearStem, tc.month, tc.hour)
		if s != tc.wantStem || b != tc.wantBranch {
			t.Errorf("%s: got %d/%d, want %d/%d", tc.name, s, b, tc.wantStem, tc.wantBranch)
		}
	}
}

// The life palace must always be a valid pair: same parity of stem and branch
// index, so that a Nayin exists for it.
func TestLifePalaceIsAlwaysAValidPair(t *testing.T) {
	for year := 1; year <= 10; year++ {
		for month := 1; month <= 12; month++ {
			for hour := 1; hour <= 12; hour++ {
				s, b := lifePalace(year, month, hour)
				if g := utils.GanZhiOf(utils.CalculateHeavenlyStem(s), utils.CalculateEarthlyBranch(b)); g.Name == "" {
					t.Fatalf("lifePalace(%d, %d, %d) = %d/%d has no Nayin", year, month, hour, s, b)
				}
			}
		}
	}
}

func TestAuxiliaryPillarsAreValidForEveryMonthAndDayPillar(t *testing.T) {
	for k := 0; k < 60; k++ {
		stem, branch := k%10+1, branchValue(k%12)
		for name, f := range map[string]func(int, int) (int, int){"thaiNguyen": thaiNguyen, "thaiTuc": thaiTuc} {
			s, b := f(stem, branch)
			if g := utils.GanZhiOf(utils.CalculateHeavenlyStem(s), utils.CalculateEarthlyBranch(b)); g.Name == "" {
				t.Fatalf("%s(%d, %d) = %d/%d has no Nayin", name, stem, branch, s, b)
			}
		}
	}
}
