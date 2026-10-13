package analysis

import (
	"testing"

	"github.com/tommitoan/bazica/v2/model"
)

func TestPatternCandidatesDistinguishExposedAndUnexposedMonthQi(t *testing.T) {
	chart := chartOf(chart1995)
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	candidates := chart.Analysis.PatternCandidates
	if len(candidates) != 1 || candidates[0].Basis != "month_main_qi_unexposed" ||
		candidates[0].HiddenStem.Value != 4 || candidates[0].TenGod.Code != "ten_god.direct_officer" ||
		len(candidates[0].ExposedAt) != 0 {
		t.Errorf("Wu month without exposed Ding: %+v", candidates)
	}

	// Yin month hides Jia, Bing and Wu; all three are exposed outside the Day Master.
	chart = chartOf([4][2]int{{1, 1}, {3, 1}, {7, 5}, {5, 1}})
	if err := Attach(chart, model.GenderMale); err != nil {
		t.Fatal(err)
	}
	candidates = chart.Analysis.PatternCandidates
	if len(candidates) != 3 {
		t.Fatalf("exposed Yin month candidates = %+v", candidates)
	}
	for i, want := range []struct {
		rank   int
		pillar string
	}{{1, "year"}, {2, "month"}, {3, "hour"}} {
		if candidates[i].QiRank != want.rank || len(candidates[i].ExposedAt) != 1 || candidates[i].ExposedAt[0] != want.pillar ||
			candidates[i].Basis != "month_hidden_stem_exposed" {
			t.Errorf("candidate %d = %+v", i, candidates[i])
		}
	}
}
