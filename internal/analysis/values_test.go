package analysis

import "testing"

func TestBranchNumberingRoundTrips(t *testing.T) {
	// Public values run Tiger = 1 to Ox = 12; indices run Rat = 0 to Pig = 11.
	for _, tc := range []struct{ value, index int }{{1, 2}, {5, 6}, {10, 11}, {11, 0}, {12, 1}} {
		if got := branchIndex(tc.value); got != tc.index {
			t.Errorf("branchIndex(%d) = %d, want %d", tc.value, got, tc.index)
		}
		if got := branchValue(tc.index); got != tc.value {
			t.Errorf("branchValue(%d) = %d, want %d", tc.index, got, tc.value)
		}
	}
	for v := 1; v <= 12; v++ {
		if got := branchValue(branchIndex(v)); got != v {
			t.Errorf("round trip of branch %d = %d", v, got)
		}
	}
}

func TestIndexHelpersWrapNegativeAndLargeInputs(t *testing.T) {
	if got := stemValue(-1); got != 10 {
		t.Errorf("stemValue(-1) = %d, want 10", got)
	}
	if got := stemValue(10); got != 1 {
		t.Errorf("stemValue(10) = %d, want 1", got)
	}
	if got := branchValue(-1); got != 10 {
		t.Errorf("branchValue(-1) = %d, want 10 (Pig)", got)
	}
}
