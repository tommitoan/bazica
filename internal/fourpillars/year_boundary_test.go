package fourpillars

import (
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

// sexagenaryIndex returns the position (0-59) of a stem and branch in the sixty-cycle.
func sexagenaryIndex(t *testing.T, stem model.HeavenlyStem, branch model.EarthlyBranch) int {
	t.Helper()
	// Branch values run from 1 (Tiger) to 12 (Ox); the classical order starts at the Rat.
	zhi := (branch.Value + 1) % 12
	gan := stem.Value - 1
	for n := 0; n < 60; n++ {
		if n%10 == gan && n%12 == zhi {
			return n
		}
	}
	t.Fatalf("no sexagenary position for stem %d branch %d", stem.Value, branch.Value)
	return 0
}

// fiveTigersStem is the stem of the Tiger month in a year with the given year stem
// (Jia and Ji years open with Bing, Yi and Geng with Wu, and so on), advanced by
// the number of months since the Tiger month.
func fiveTigersStem(yearStem, monthsSinceTiger int) int {
	tiger := ((yearStem-1)%5)*2 + 3
	return (tiger-1+monthsSinceTiger)%10 + 1
}

// The year pillar changes exactly at each Lichun and nowhere else, the month
// pillar advances one sixty-cycle step per month, and the month stem follows the
// Five Tigers rule from the year stem, for every instant sampled in the whole range.
func TestYearAndMonthPillarsAreConsistentOverTheWholeRange(t *testing.T) {
	first, last, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	step := 6 * time.Hour
	start := time.Date(first, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(last, 12, 31, 18, 0, 0, 0, time.UTC)

	var prev *model.FourPillars
	var prevAt time.Time
	yearChanges := 0
	for at := start; !at.After(end); at = at.Add(step) {
		p, _, _, err := GetFourPillars(at, time.UTC)
		if err != nil {
			t.Fatalf("%s: %v", at, err)
		}
		months := (p.MonthPillar.EarthlyBranch.Value - 1)
		if want, got := fiveTigersStem(p.YearPillar.HeavenlyStem.Value, months), p.MonthPillar.HeavenlyStem.Value; want != got {
			t.Fatalf("%s: month stem %d, Five Tigers rule gives %d (year stem %d, month branch %d)",
				at, got, want, p.YearPillar.HeavenlyStem.Value, p.MonthPillar.EarthlyBranch.Value)
		}
		if prev != nil {
			yearNow := sexagenaryIndex(t, p.YearPillar.HeavenlyStem, p.YearPillar.EarthlyBranch)
			yearBefore := sexagenaryIndex(t, prev.YearPillar.HeavenlyStem, prev.YearPillar.EarthlyBranch)
			monthNow := sexagenaryIndex(t, p.MonthPillar.HeavenlyStem, p.MonthPillar.EarthlyBranch)
			monthBefore := sexagenaryIndex(t, prev.MonthPillar.HeavenlyStem, prev.MonthPillar.EarthlyBranch)

			lichun, err := utils.StartOfSpring(at.Year())
			if err != nil {
				t.Fatal(err)
			}
			lichunBetween := lichun.After(prevAt) && !lichun.After(at)

			switch {
			case yearNow != yearBefore:
				yearChanges++
				if !lichunBetween {
					t.Fatalf("%s: the year pillar changed without a Lichun between %s and it", at, prevAt)
				}
				if yearNow != (yearBefore+1)%60 {
					t.Fatalf("%s: the year pillar moved from %d to %d, not one step", at, yearBefore, yearNow)
				}
				if p.MonthPillar.EarthlyBranch.Value != 1 {
					t.Fatalf("%s: the year changed but the month branch is %d, not the Tiger", at, p.MonthPillar.EarthlyBranch.Value)
				}
			case lichunBetween:
				t.Fatalf("%s: a Lichun lies between %s and it but the year pillar did not change", at, prevAt)
			}
			if monthNow != monthBefore && monthNow != (monthBefore+1)%60 {
				t.Fatalf("%s: the month pillar moved from %d to %d, not one step", at, monthBefore, monthNow)
			}
		}
		prev, prevAt = p, at
	}
	// One change per Lichun between the first and last sample.
	if want := last - first + 1; yearChanges < want-1 || yearChanges > want {
		t.Errorf("the year pillar changed %d times over %d years", yearChanges, want)
	}
}

// One minute before Lichun still belongs to the previous Ba-zi year and the
// instant of Lichun starts the new one, whatever the birth zone, including the
// late hours that move the day pillar.
func TestYearChangesAtTheLichunInstantInEveryZone(t *testing.T) {
	first, last, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	for _, hours := range []int{-12, -5, 0, 7, 8, 14} {
		zone := time.FixedZone("test", hours*3600)
		for year := first; year <= last; year++ {
			lichun, err := utils.StartOfSpring(year)
			if err != nil {
				t.Fatal(err)
			}
			before, err := GetBaziYear(lichun.Add(-time.Minute).In(zone))
			if err != nil {
				t.Fatal(err)
			}
			at, err := GetBaziYear(lichun.In(zone))
			if err != nil {
				t.Fatal(err)
			}
			if before != year-1 || at != year {
				t.Fatalf("zone %+d, Lichun %s: bazi year %d one minute before, %d at the instant; want %d and %d",
					hours, lichun, before, at, year-1, year)
			}
		}
	}
}

// The ends of the supported range: the first day belongs to the Ba-zi year before it
// (it precedes Lichun) and the last day to its own year.
func TestBaziYearAtBothEndsOfTheRange(t *testing.T) {
	first, last, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	if y, err := GetBaziYear(time.Date(first, 1, 1, 0, 0, 0, 0, time.UTC)); err != nil || y != first-1 {
		t.Errorf("first day: year %d, error %v; want %d", y, err, first-1)
	}
	if y, err := GetBaziYear(time.Date(last, 12, 31, 23, 30, 0, 0, time.UTC)); err != nil || y != last {
		t.Errorf("last day: year %d, error %v; want %d", y, err, last)
	}
}
