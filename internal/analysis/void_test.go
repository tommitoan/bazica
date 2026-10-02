package analysis

import "testing"

// For each of the sixty day pillars, the void branches are the two branches
// that do not occur among the ten pillars of its cycle.
func TestVoidBranchesForAllSixtyDays(t *testing.T) {
	for k := 0; k < 60; k++ {
		stem, branch := k%10+1, branchValue(k%12)
		present := map[int]bool{}
		start := k - k%10
		for j := 0; j < 10; j++ {
			present[branchValue((start+j)%12)] = true
		}
		var missing []int
		for b := 1; b <= 12; b++ {
			if !present[b] {
				missing = append(missing, b)
			}
		}
		got := voidBranches(stem, branch)
		if len(missing) != 2 || !(got[0] == missing[0] && got[1] == missing[1] || got[0] == missing[1] && got[1] == missing[0]) {
			t.Errorf("day %d (stem %d, branch %d): got %v, want %v", k, stem, branch, got, missing)
		}
	}
}

func TestVoidBranchesKnownDays(t *testing.T) {
	// Jia-Zi: Xu and Hai are void. Geng-Wu (1995 chart): the same pair.
	for _, tc := range []struct{ stem, branch int }{{1, 11}, {7, 5}} {
		if got := voidBranches(tc.stem, tc.branch); got != [2]int{9, 10} {
			t.Errorf("voidBranches(%d, %d) = %v, want [9 10]", tc.stem, tc.branch, got)
		}
	}
}
