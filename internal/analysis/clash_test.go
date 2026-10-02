package analysis

import "testing"

func TestOvercomesIsDirectionalAndNeedsSamePolarity(t *testing.T) {
	for _, tc := range []struct {
		stem, other int
		want        bool
	}{
		{3, 7, true},  // Bing over Geng
		{7, 3, false}, // Geng does not overcome Bing
		{3, 8, false}, // Bing over Xin: different polarity
		{5, 9, true},  // Wu over Ren: earth takes part
		{1, 5, true},  // Jia over Wu
		{7, 1, true},  // Geng over Jia
		{4, 8, true},  // Ding over Xin
		{1, 1, false}, // never against itself
	} {
		if got := overcomes(tc.stem, tc.other); got != tc.want {
			t.Errorf("overcomes(%d, %d) = %v, want %v", tc.stem, tc.other, got, tc.want)
		}
	}
}

func TestOvercomesPairCount(t *testing.T) {
	n := 0
	for a := 1; a <= 10; a++ {
		for b := 1; b <= 10; b++ {
			if overcomes(a, b) {
				n++
			}
		}
	}
	if n != 10 {
		t.Errorf("overcoming pairs = %d, want 10", n)
	}
}

func TestHeavenEarthClash(t *testing.T) {
	// Bing-Zi (3, Rat 11) against Geng-Wu (7, Horse 5): stem overcomes and branches are opposite.
	bingZi, gengWu := cell{3, 11}, cell{7, 5}
	if !heavenEarthClash(bingZi, []cell{gengWu}) {
		t.Error("Bing-Zi should be flagged against Geng-Wu")
	}
	if heavenEarthClash(gengWu, []cell{bingZi}) {
		t.Error("Geng-Wu must not be flagged by the same pair")
	}
	if heavenEarthClash(bingZi, []cell{{7, 6}}) {
		t.Error("branches that are not opposite must not clash")
	}
	if heavenEarthClash(bingZi, nil) {
		t.Error("no counterparts, no clash")
	}
}
