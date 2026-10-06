package utils

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	dataTables "github.com/tommitoan/bazica/v2/data"
	"github.com/tommitoan/bazica/v2/model"
)

// The solar-term table was regenerated from the JPL DE440 ephemeris. These tests
// compare it with the table of v1.4.3 (testdata/calendar) to show exactly what
// moved inside the years that release covered.

func loadFixtureSolarTerms(t *testing.T, path string) map[int][]solarTermTime {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var table map[string]model.CombinedData
	if err := json.Unmarshal(raw, &table); err != nil {
		t.Fatal(err)
	}
	years := make(map[int][]solarTermTime, len(table))
	for key, entry := range table {
		year, err := strconv.Atoi(key)
		if err != nil {
			t.Fatalf("year key %q", key)
		}
		terms, err := parseSolarTerms(entry.Data)
		if err != nil {
			t.Fatal(err)
		}
		years[year] = terms
	}
	return years
}

// Every term of 1900-2100 moves by less than 100 seconds, which is the Delta T
// model difference between the two tables, except the 1903 winter solstice that
// the old table had about twenty minutes early. The margin year 1899 is not
// compared: the old table was off by up to hours there for its first 22 terms.
func TestSolarTermsMovedLittleSinceV143(t *testing.T) {
	old := loadFixtureSolarTerms(t, "../../testdata/calendar/solar-term-v1.4.3.json")
	current, err := loadSolarTerms()
	if err != nil {
		t.Fatal(err)
	}
	var largest time.Duration
	for year := 1900; year <= 2100; year++ {
		for i, was := range old[year] {
			now := current[year][i]
			if now.name != was.name {
				t.Fatalf("%d: term %d is %s, was %s", year, i, now.name, was.name)
			}
			diff := now.time.Sub(was.time).Abs()
			if year == 1903 && was.name == model.WinterSolstice {
				if diff < 19*time.Minute || diff > 21*time.Minute {
					t.Errorf("1903 winter solstice moved by %s, want about 20 minutes", diff)
				}
				continue
			}
			if diff > 100*time.Second {
				t.Errorf("%d %s moved by %s", year, was.name, diff)
			}
			if diff > largest {
				largest = diff
			}
		}
	}
	t.Logf("largest movement outside the 1903 winter solstice: %s", largest)
}

// A birth that is not within two minutes of a jie keeps its month pillar and its
// solar-term minutes within that margin, so the regenerated table only affects
// births on the very edge of a month.
func TestRegeneratedTableKeepsTheMonthOfBirthsAwayFromTheEdge(t *testing.T) {
	old := loadFixtureSolarTerms(t, "../../testdata/calendar/solar-term-v1.4.3.json")
	for year := 1900; year <= 2099; year++ {
		for _, was := range old[year] {
			if !isInitial(was.name) || (year == 1903 && was.name == model.WinterSolstice) {
				continue
			}
			margin := 2 * time.Minute
			for _, probe := range []time.Time{was.time.Add(margin), was.time.Add(-margin)} {
				gotNew, _, _, errNew := GetSolarTerm(probe)
				gotOld, _, _, errOld := findSolarTermFor(probe, old)
				if errNew != nil || errOld != nil {
					t.Fatalf("%d %s: %v %v", year, was.name, errNew, errOld)
				}
				if gotNew != gotOld {
					t.Errorf("%d %s: probe %s is %s now, was %s", year, was.name, probe, gotNew, gotOld)
				}
			}
		}
	}
}

func isInitial(name string) bool {
	for _, n := range InitialTerms {
		if n == name {
			return true
		}
	}
	return false
}

func findSolarTermFor(t time.Time, years map[int][]solarTermTime) (string, int, int, error) {
	year := t.Year()
	return findSolarTerm(t, years[year], years[year+1], years[year-1])
}

// Lunar New Year of 1900-2099 is unchanged: the generator reproduced all 200 dates.
func TestLunarNewYearUnchangedSinceV143(t *testing.T) {
	raw, err := os.ReadFile("../../testdata/calendar/lunar-new-year-v1.4.3.json")
	if err != nil {
		t.Fatal(err)
	}
	var old, current model.LunarNewYearData
	if err := json.Unmarshal(raw, &old); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(dataTables.LunarNewYear, &current); err != nil {
		t.Fatal(err)
	}
	if len(old.LunarNewYearDates) != 200 {
		t.Fatalf("fixture has %d years, want 200", len(old.LunarNewYearDates))
	}
	for year, want := range old.LunarNewYearDates {
		if got := current.LunarNewYearDates[year]; got != want {
			t.Errorf("%s: Lunar New Year is %s, was %s", year, got, want)
		}
	}
}
