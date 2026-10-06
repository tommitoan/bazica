package analysis

import (
	"slices"
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/internal/fourpillars"
)

// For a birth between the Lunar New Year and Lichun the stars differ from what the
// Lunar New Year year would have given: many rules start from the year pillar and most
// stars can land in the year or month pillar, and both pillars change. The README tells
// readers to treat the stars of these births as new, so the statement is pinned here.
func TestStarsChangeForEveryBirthBetweenTheLunarNewYearAndLichun(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*3600)
	classicalBranch := func(branchValue int) int { return (branchValue + 1) % 12 } // Rat is 0
	terms := map[string]bool{}
	charts := 0

	for day := time.Date(1901, 1, 2, 12, 0, 0, 0, loc); day.Year() < 2099; day = day.AddDate(0, 0, 1) {
		baziYear, err := fourpillars.GetBaziYear(day)
		if err != nil {
			t.Fatal(err)
		}
		lunarYear, err := fourpillars.GetLunarYear(day)
		if err != nil {
			t.Fatal(err)
		}
		if baziYear == lunarYear {
			continue
		}
		p, _, _, err := fourpillars.GetFourPillars(day, loc)
		if err != nil {
			t.Fatal(err)
		}
		stems := [4]int{p.YearPillar.HeavenlyStem.Value - 1, p.MonthPillar.HeavenlyStem.Value - 1, p.DayPillar.HeavenlyStem.Value - 1, p.HourPillar.HeavenlyStem.Value - 1}
		branches := [4]int{classicalBranch(p.YearPillar.EarthlyBranch.Value), classicalBranch(p.MonthPillar.EarthlyBranch.Value),
			classicalBranch(p.DayPillar.EarthlyBranch.Value), classicalBranch(p.HourPillar.EarthlyBranch.Value)}

		// The pillars the Lunar New Year rule would have produced: its year, and the month stem
		// that follows from that year by the Five Tigers rule (the branch is the same).
		lunarStem := ((lunarYear-4)%10 + 10) % 10
		lunarStems, lunarBranches := stems, branches
		lunarStems[0], lunarBranches[0] = lunarStem, ((lunarYear-4)%12+12)%12
		tigerMonthStem := (lunarStem%5)*2 + 2
		lunarStems[1] = (tigerMonthStem + p.MonthPillar.EarthlyBranch.Value - 1) % 10

		for _, male := range []bool{false, true} {
			charts++
			now := evaluateStars(facts(stems, branches, male))
			before := evaluateStars(facts(lunarStems, lunarBranches, male))
			differs := false
			for pillar := range now {
				codes := func(list []string) []string { slices.Sort(list); return list }
				var a, b []string
				for _, s := range now[pillar] {
					a = append(a, s.Code)
				}
				for _, s := range before[pillar] {
					b = append(b, s.Code)
				}
				if !slices.Equal(codes(a), codes(b)) {
					differs = true
				}
				for _, code := range a {
					if !slices.Contains(b, code) {
						terms[code] = true
					}
				}
				for _, code := range b {
					if !slices.Contains(a, code) {
						terms[code] = true
					}
				}
			}
			if !differs {
				t.Fatalf("%s (male %v): the stars equal those of the Lunar New Year pillars", day.Format("2006-01-02"), male)
			}
		}
	}
	if charts < 2800 {
		t.Fatalf("only %d charts examined, expected about 2900", charts)
	}
	if len(terms) < 30 {
		t.Errorf("%d star terms changed, expected about 39", len(terms))
	}
}
