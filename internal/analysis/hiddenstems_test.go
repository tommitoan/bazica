package analysis

import (
	"reflect"
	"testing"
)

func TestHiddenStemsByBranch(t *testing.T) {
	// Branch values: Tiger 1 ... Pig 10, Rat 11, Ox 12. Main qi first.
	want := map[int][]int{
		1: {1, 3, 5}, 2: {2}, 3: {5, 2, 10}, 4: {3, 5, 7}, 5: {4, 6}, 6: {6, 4, 2},
		7: {7, 9, 5}, 8: {8}, 9: {5, 8, 4}, 10: {9, 1}, 11: {10}, 12: {6, 10, 8},
	}
	for branch, stems := range want {
		if got := hiddenStemsOf(branch); !reflect.DeepEqual(got, stems) {
			t.Errorf("hiddenStemsOf(%d) = %v, want %v", branch, got, stems)
		}
	}
}

func TestHiddenStemsAreValidStems(t *testing.T) {
	total := 0
	for b := 1; b <= 12; b++ {
		for _, s := range hiddenStemsOf(b) {
			if !validStem(s) {
				t.Errorf("branch %d hides invalid stem %d", b, s)
			}
			total++
		}
	}
	if total != 28 {
		t.Errorf("total hidden stems = %d, want 28", total)
	}
}
