package fourpillars

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
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
	stepDays := 1
	if testing.Short() {
		stepDays = 3
	}
	for day := time.Date(first, 1, 2, 12, 0, 0, 0, loc); day.Year() < last; day = day.AddDate(0, 0, stepDays) {
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
	if want := 5200 / stepDays; differing < want*85/100 || differing > want*115/100 {
		t.Errorf("%d of %d days differ, expected about %d", differing, days, want)
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

// The lunar calendar changes year at midnight. The 23:00 rule moves the Ba-zi day, not the
// lunar year, so the eve of the Lunar New Year 1990 (27 January) still belongs to 1989
// until 00:00.
func TestLunarYearChangesAtMidnightNotAtTheRatHour(t *testing.T) {
	loc := time.FixedZone("UTC+8", 8*3600)
	tests := []struct {
		at   time.Time
		want int
	}{
		{time.Date(1990, 1, 26, 22, 59, 0, 0, loc), 1989},
		{time.Date(1990, 1, 26, 23, 0, 0, 0, loc), 1989},
		{time.Date(1990, 1, 26, 23, 59, 0, 0, loc), 1989},
		{time.Date(1990, 1, 27, 0, 0, 0, 0, loc), 1990},
		{time.Date(1990, 1, 27, 23, 30, 0, 0, loc), 1990},
	}
	for _, tc := range tests {
		got, err := GetLunarYear(tc.at)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("%s: lunar year %d, want %d", tc.at.Format("2006-01-02 15:04"), got, tc.want)
		}
	}
}

// Dates outside the tables fail with ErrDateOutOfRange instead of returning a made-up
// year. GetBaziYear needs the Lichun of one year; GetYearReference also needs the Lunar
// New Year table, whose range is the supported birth range.
func TestYearReferenceOutsideTheTables(t *testing.T) {
	first, last, err := SupportedYears()
	if err != nil {
		t.Fatal(err)
	}
	for _, year := range []int{1500, 3000} {
		at := time.Date(year, 6, 1, 12, 0, 0, 0, time.UTC)
		if _, err := GetBaziYear(at); !errors.Is(err, model.ErrDateOutOfRange) {
			t.Errorf("GetBaziYear(%d) error = %v, want ErrDateOutOfRange", year, err)
		}
	}
	for _, year := range []int{1500, first - 1, last + 1, 3000} {
		at := time.Date(year, 6, 1, 12, 0, 0, 0, time.UTC)
		if _, err := GetYearReference(at, &model.YearPillar{}); !errors.Is(err, model.ErrDateOutOfRange) {
			t.Errorf("GetYearReference(%d) error = %v, want ErrDateOutOfRange", year, err)
		}
	}
}
