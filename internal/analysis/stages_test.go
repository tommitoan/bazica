package analysis

import "testing"

// birthBranch is the branch value where each stem is at its Birth stage.
var birthBranch = [10]int{10, 5, 1, 8, 1, 8, 4, 11, 7, 2}

var stageCodes = []string{
	"birth", "bath", "youth", "thriving", "prosperous", "weakening",
	"sick", "death", "grave", "extinction", "conceived", "nourishing",
}

// Yang stems advance along the branches from their Birth branch; Yin stems run against them.
func TestStageForEveryStemAndBranch(t *testing.T) {
	for stem := 1; stem <= 10; stem++ {
		for branch := 1; branch <= 12; branch++ {
			step := branch - birthBranch[stem-1]
			if stem%2 == 0 {
				step = -step
			}
			want := "stage." + stageCodes[((step%12)+12)%12]
			if got := stage(stem, branch).Code; got != want {
				t.Errorf("stem %d at branch %d: got %s, want %s", stem, branch, got, want)
			}
		}
	}
}

func TestStageKnownChart(t *testing.T) {
	// Geng at Hai (Pig) is Sick in the 1995 reference chart.
	if got := stage(7, 10); got.Code != "stage.sick" || got.VI != "Bệnh" {
		t.Errorf("stage(7, 10) = %+v", got)
	}
}
