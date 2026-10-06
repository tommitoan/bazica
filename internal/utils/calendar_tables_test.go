package utils

import (
	"testing"

	"github.com/tommitoan/bazica/v2/model"
)

func TestCalendarTablesAreIndexedByValue(t *testing.T) {
	for i, s := range heavenlyStems {
		if s.Value != i+1 || s.Name == "" || s.Character == "" || s.Spelling == "" {
			t.Errorf("heavenlyStems[%d] = %+v", i, s)
		}
	}
	for i, b := range earthlyBranches {
		if b.Value != i+1 || b.Name == "" || b.Character == "" || b.Spelling == "" {
			t.Errorf("earthlyBranches[%d] = %+v", i, b)
		}
	}
	for i, lc := range lifeCycles {
		if lc == "" {
			t.Errorf("lifeCycles[%d] is empty", i)
		}
	}
}

func TestEveryTermMapsToABranch(t *testing.T) {
	terms := append(append([]string{}, InitialTerms...), MidpointTerms...)
	if len(terms) != 24 || len(termBranchValue) != 24 {
		t.Fatalf("got %d terms and %d mappings, want 24 each", len(terms), len(termBranchValue))
	}
	for _, term := range terms {
		if ConvertTermToBranch(term).Value == 0 {
			t.Errorf("term %s has no branch", term)
		}
	}
	// Initial and midpoint terms of one month share a branch.
	for i := range InitialTerms {
		if a, b := ConvertTermToBranch(InitialTerms[i]), ConvertTermToBranch(MidpointTerms[i]); a != b {
			t.Errorf("%s and %s map to different branches", InitialTerms[i], MidpointTerms[i])
		}
	}
}

func TestOutOfRangeValuesYieldZeroValues(t *testing.T) {
	// Stems and the rules keyed by stem run 1-10; branches and life stages 1-12.
	for _, n := range []int{-1, 0, 11, 13} {
		if got := CalculateHeavenlyStem(n); got != (model.HeavenlyStem{}) {
			t.Errorf("CalculateHeavenlyStem(%d) = %+v, want zero", n, got)
		}
		if got := GetLifeCycleRule(n); got != 0 {
			t.Errorf("GetLifeCycleRule(%d) = %d, want 0", n, got)
		}
		if GetStemRuleByFiveTigers(n) != 0 || GetStemRuleByFiveRats(n) != 0 {
			t.Errorf("five-tigers/rats rule for %d should be 0", n)
		}
	}
	for _, n := range []int{-1, 0, 13} {
		if got := CalculateEarthlyBranch(n); got != (model.EarthlyBranch{}) {
			t.Errorf("CalculateEarthlyBranch(%d) = %+v, want zero", n, got)
		}
		if got := ConvertValueToLifeCycle(n); got != "" {
			t.Errorf("ConvertValueToLifeCycle(%d) = %q, want empty", n, got)
		}
	}
	if got := ConvertTermToBranch("unknown"); got != (model.EarthlyBranch{}) {
		t.Errorf("ConvertTermToBranch(unknown) = %+v, want zero", got)
	}
}

func TestStemRules(t *testing.T) {
	// Jia/Ji years begin with Bing Yin; Jia/Ji days begin with Jia Zi.
	tigers := map[int]int{1: 3, 6: 3, 2: 5, 7: 5, 3: 7, 8: 7, 4: 9, 9: 9, 5: 1, 10: 1}
	rats := map[int]int{1: 1, 6: 1, 2: 3, 7: 3, 3: 5, 8: 5, 4: 7, 9: 7, 5: 9, 10: 9}
	for stem := 1; stem <= 10; stem++ {
		if got := GetStemRuleByFiveTigers(stem); got != tigers[stem] {
			t.Errorf("GetStemRuleByFiveTigers(%d) = %d, want %d", stem, got, tigers[stem])
		}
		if got := GetStemRuleByFiveRats(stem); got != rats[stem] {
			t.Errorf("GetStemRuleByFiveRats(%d) = %d, want %d", stem, got, rats[stem])
		}
	}
}

// Stems and branches pair up only when they share polarity, which yields the
// 60 combinations of the sexagenary cycle.
func TestGanZhiOfCoversTheSexagenaryCycle(t *testing.T) {
	pairs := 0
	for _, stem := range heavenlyStems {
		for _, branch := range earthlyBranches {
			g := GanZhiOf(stem, branch)
			if stem.Value%2 != branch.Value%2 {
				if g.Name != "" {
					t.Errorf("%s %s should not be a valid pair", stem.Name, branch.Name)
				}
				continue
			}
			pairs++
			if g.Name == "" || g.ElementName == "" || g.ElementValue == 0 {
				t.Errorf("%s %s has incomplete data: %+v", stem.Name, branch.Name, g)
			}
		}
	}
	if pairs != 60 {
		t.Errorf("got %d valid pairs, want 60", pairs)
	}
}

func TestGetGanzhiUnknownSign(t *testing.T) {
	if name, element, value := GetGanzhi("bogus"); name != "" || element != "" || value != 0 {
		t.Errorf("GetGanzhi(bogus) = %q %q %d", name, element, value)
	}
}
