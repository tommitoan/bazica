package utils

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	dataTables "github.com/tommitoan/bazica/data"
	"github.com/tommitoan/bazica/model"
)

// Guards the bundled calendar data: every year must list the 24 terms in
// strictly increasing order and each term must fall in its own month.
func TestSolarTermDataIntegrity(t *testing.T) {
	years, err := loadSolarTerms()
	if err != nil {
		t.Fatal(err)
	}

	months := []time.Month{
		time.January, time.January, time.February, time.February, time.March, time.March,
		time.April, time.April, time.May, time.May, time.June, time.June,
		time.July, time.July, time.August, time.August, time.September, time.September,
		time.October, time.October, time.November, time.November, time.December, time.December,
	}
	for year := 1899; year <= 2100; year++ {
		terms, ok := years[year]
		if !ok {
			t.Errorf("%d: missing", year)
			continue
		}
		if len(terms) != 24 {
			t.Errorf("%d: has %d terms, want 24", year, len(terms))
			continue
		}
		for i, term := range terms {
			u := term.time.UTC()
			if u.Year() != year || u.Month() != months[i] {
				t.Errorf("%d: %s is dated %s", year, term.name, u.Format("2006-01-02"))
			}
			if i > 0 && !term.time.After(terms[i-1].time) {
				t.Errorf("%d: %s is not after %s", year, term.name, terms[i-1].name)
			}
		}
	}
}

// Lunar years last 12 or 13 lunations, which is 353-355 or 383-385 days.
func TestLunarNewYearGaps(t *testing.T) {
	var data model.LunarNewYearData
	if err := json.Unmarshal(dataTables.LunarNewYear, &data); err != nil {
		t.Fatal(err)
	}
	dates := map[int]time.Time{}
	for year := 1900; year <= 2099; year++ {
		d, err := time.Parse("2006-01-02", strconv.Itoa(year)+"-"+data.LunarNewYearDates[strconv.Itoa(year)])
		if err != nil {
			t.Fatalf("%d: %v", year, err)
		}
		dates[year] = d
	}
	for year := 1900; year < 2099; year++ {
		gap := int(dates[year+1].Sub(dates[year]).Hours() / 24)
		if !(gap >= 353 && gap <= 355) && !(gap >= 383 && gap <= 385) {
			t.Errorf("%d to %d: %d days between Lunar New Years", year, year+1, gap)
		}
	}
}

// A handful of widely published dates, to catch a table keyed by the wrong year.
func TestLunarNewYearKnownDates(t *testing.T) {
	var data model.LunarNewYearData
	if err := json.Unmarshal(dataTables.LunarNewYear, &data); err != nil {
		t.Fatal(err)
	}
	known := map[int]string{
		1900: "01-31", 1949: "01-29", 1950: "02-17", 1958: "02-18", 1985: "02-20",
		1990: "01-27", 1995: "01-31", 2000: "02-05", 2020: "01-25", 2034: "02-19",
	}
	for year, want := range known {
		if got := data.LunarNewYearDates[strconv.Itoa(year)]; got != want {
			t.Errorf("%d: Lunar New Year is %s, want %s", year, got, want)
		}
	}
}

func TestLunarNewYearDataIntegrity(t *testing.T) {
	var data model.LunarNewYearData
	if err := json.Unmarshal(dataTables.LunarNewYear, &data); err != nil {
		t.Fatal(err)
	}
	for year := 1900; year <= 2099; year++ {
		v, ok := data.LunarNewYearDates[strconv.Itoa(year)]
		if !ok {
			t.Errorf("%d: missing", year)
			continue
		}
		d, err := time.Parse("2006-01-02", strconv.Itoa(year)+"-"+v)
		if err != nil {
			t.Errorf("%d: bad date %q", year, v)
			continue
		}
		if d.YearDay() < 21 || d.YearDay() > 51 { // Jan 21 - Feb 20
			t.Errorf("%d: Lunar New Year %s is outside Jan 21 - Feb 20", year, v)
		}
	}
}
