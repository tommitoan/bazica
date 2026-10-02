package ultis

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/tommitoan/bazica/model"
)

// Guards the bundled calendar data: every year must list the 24 terms in
// strictly increasing order and each term must fall in its own month.
func TestSolarTermDataIntegrity(t *testing.T) {
	raw, err := os.ReadFile(dataPath + "data/solar-term.json")
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]model.CombinedData
	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}

	months := []time.Month{
		time.January, time.January, time.February, time.February, time.March, time.March,
		time.April, time.April, time.May, time.May, time.June, time.June,
		time.July, time.July, time.August, time.August, time.September, time.September,
		time.October, time.October, time.November, time.November, time.December, time.December,
	}
	for year := 1899; year <= 2100; year++ {
		entry, ok := data[strconv.Itoa(year)]
		if !ok {
			t.Errorf("%d: missing", year)
			continue
		}
		terms, err := parseSolarTerms(entry.Data)
		if err != nil {
			t.Errorf("%d: %v", year, err)
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

func TestLunarNewYearDataIntegrity(t *testing.T) {
	raw, err := os.ReadFile(dataPath + "data/lunar-new-year.json")
	if err != nil {
		t.Fatal(err)
	}
	var data model.LunarNewYearData
	if err := json.Unmarshal(raw, &data); err != nil {
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
