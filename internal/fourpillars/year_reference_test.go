package fourpillars

import (
	"fmt"
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/internal/utils"
)

// Differs is true exactly when the Lunar New Year and Lichun fall on different sides of
// the day, for every day of the range at noon UTC+8. The expectation is computed from
// the two tables with plain date arithmetic, not from the code under test.
func TestYearReferenceDiffersExactlyBetweenTheLunarNewYearAndLichun(t *testing.T) {
	first, last, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	loc := time.FixedZone("UTC+8", 8*3600)
	lunarData, err := loadLunarNewYear()
	if err != nil {
		t.Fatal(err)
	}

	days, differing := 0, 0
	for day := time.Date(first, 1, 2, 12, 0, 0, 0, loc); day.Year() < last; day = day.AddDate(0, 0, 1) {
		year := day.Year()
		lichun, err := utils.StartOfSpring(year)
		if err != nil {
			t.Fatal(err)
		}
		newYear, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%d-%s", year, lunarData.LunarNewYearDates[fmt.Sprint(year)]), loc)
		if err != nil {
			t.Fatal(err)
		}
		afterLichun := !day.Before(lichun)
		afterNewYear := !day.Before(newYear)
		want := afterLichun != afterNewYear

		pillars, _, _, err := GetFourPillars(day, loc)
		if err != nil {
			t.Fatal(err)
		}
		reference, err := GetYearReference(day, pillars.YearPillar)
		if err != nil {
			t.Fatal(err)
		}
		if reference.Differs != want {
			t.Fatalf("%s: Differs = %v, want %v (Lichun %s, Lunar New Year %s)", day.Format("2006-01-02"), reference.Differs, want, lichun.In(loc), newYear.Format("2006-01-02"))
		}
		days++
		if want {
			differing++
		}
	}
	// About 7.4 days a year over 700 years.
	if differing < 4500 || differing > 6000 {
		t.Errorf("%d of %d days differ, expected about 5200", differing, days)
	}
}

// The reference of a birth between the Lunar New Year and Lichun names both years.
func TestYearReferenceNamesBothYears(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*3600)
	birth := time.Date(1990, 1, 31, 12, 0, 0, 0, loc)
	pillars, _, _, err := GetFourPillars(birth, loc)
	if err != nil {
		t.Fatal(err)
	}
	reference, err := GetYearReference(birth, pillars.YearPillar)
	if err != nil {
		t.Fatal(err)
	}
	if got := pillars.YearPillar.HeavenlyStem.Character + pillars.YearPillar.EarthlyBranch.Character; got != "己巳" {
		t.Errorf("year pillar = %s, want 己巳", got)
	}
	lunar := reference.LunarYear
	if lunar.Year != 1990 || lunar.HeavenlyStem.Character != "庚" || lunar.EarthlyBranch.Character != "午" || lunar.EarthlyBranch.Name != "Horse" {
		t.Errorf("lunar year = %d %s%s %s, want 1990 庚午 Horse", lunar.Year, lunar.HeavenlyStem.Character, lunar.EarthlyBranch.Character, lunar.EarthlyBranch.Name)
	}
	if lunar.NewYearDate != "1990-01-27" {
		t.Errorf("new year date = %q, want 1990-01-27", lunar.NewYearDate)
	}
	if !reference.Differs || !reference.Lichun.BirthBeforeLichun {
		t.Errorf("Differs = %v, BirthBeforeLichun = %v, want both true", reference.Differs, reference.Lichun.BirthBeforeLichun)
	}
	if _, offset := reference.Lichun.Time.Zone(); offset != 8*3600 {
		t.Errorf("Lichun time is in offset %d, want the birth zone (+8h)", offset)
	}
}

// The first supported day belongs to the lunar year before the first table year, which has
// no Lunar New Year date in the table.
func TestYearReferenceOnTheFirstSupportedDay(t *testing.T) {
	first, _, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	birth := time.Date(first, 1, 1, 12, 0, 0, 0, time.UTC)
	pillars, _, _, err := GetFourPillars(birth, time.UTC)
	if err != nil {
		t.Fatal(err)
	}
	reference, err := GetYearReference(birth, pillars.YearPillar)
	if err != nil {
		t.Fatal(err)
	}
	if reference.LunarYear.Year != first-1 || reference.LunarYear.NewYearDate != "" || reference.Differs {
		t.Errorf("reference = %+v, want lunar year %d, no new year date, no difference", reference.LunarYear, first-1)
	}
}
