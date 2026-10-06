package utils

import (
	"encoding/json"
	"fmt"
	"github.com/tommitoan/bazica/v2/data"
	"github.com/tommitoan/bazica/v2/model"
	"slices"
	"strconv"
	"sync"
	"time"
)

var InitialTerms = []string{
	model.MinorCold,
	model.StartOfSpring,
	model.AwakeningOfInsects,
	model.PureBrightness,
	model.StartOfSummer,
	model.GrainInEar,
	model.MinorHeat,
	model.StartOfAutumn,
	model.WhiteDew,
	model.ColdDew,
	model.StartOfWinter,
	model.MajorSnow,
}

var MidpointTerms = []string{
	model.MajorCold,
	model.SpringShowers,
	model.SpringEquinox,
	model.GrainRain,
	model.GrainBuds,
	model.SummerSolstice,
	model.MajorHeat,
	model.EndOfHeat,
	model.AutumnEquinox,
	model.Frost,
	model.MinorSnow,
	model.WinterSolstice,
}

// GetSolarTerm returns the solar term containing dateTime together with the
// minutes elapsed since the preceding initial term and the minutes remaining
// until the next one.
func GetSolarTerm(dateTime time.Time) (string, int, int, error) {
	years, err := loadSolarTerms()
	if err != nil {
		return "", 0, 0, err
	}
	year := dateTime.Year()
	var terms [3][]solarTermTime
	for i, y := range []int{year, year + 1, year - 1} {
		t, ok := years[y]
		if !ok {
			return "", 0, 0, fmt.Errorf("%w: no solar term data for %d", model.ErrDateOutOfRange, y)
		}
		terms[i] = t
	}
	return findSolarTerm(dateTime, terms[0], terms[1], terms[2])
}

// StartOfSpring returns the instant of Lichun (Start of Spring) in the given
// Gregorian year. It is the moment the Ba-zi year changes.
func StartOfSpring(year int) (time.Time, error) {
	years, err := loadSolarTerms()
	if err != nil {
		return time.Time{}, err
	}
	terms, ok := years[year]
	if !ok {
		return time.Time{}, fmt.Errorf("%w: no solar term data for %d", model.ErrDateOutOfRange, year)
	}
	for _, term := range terms {
		if term.name == model.StartOfSpring {
			return term.time, nil
		}
	}
	return time.Time{}, fmt.Errorf("%w: no %s in the solar term data for %d", model.ErrDataUnavailable, model.StartOfSpring, year)
}

// SolarTermYears returns the first and last year the embedded solar-term table
// holds. A birth year y needs the years y-1, y and y+1.
func SolarTermYears() (first, last int, err error) {
	years, err := loadSolarTerms()
	if err != nil {
		return 0, 0, err
	}
	first, last = 0, 0
	for year := range years {
		if first == 0 || year < first {
			first = year
		}
		if year > last {
			last = year
		}
	}
	return first, last, nil
}

type solarTermTime struct {
	name string
	time time.Time
}

// parseSolarTerms returns the 24 solar terms of a year in calendar order.
func parseSolarTerms(y model.SolarTermYear) ([]solarTermTime, error) {
	raw := []struct{ name, value string }{
		{model.MinorCold, y.MinorCold},
		{model.MajorCold, y.MajorCold},
		{model.StartOfSpring, y.StartOfSpring},
		{model.SpringShowers, y.SpringShowers},
		{model.AwakeningOfInsects, y.AwakeningOfInsects},
		{model.SpringEquinox, y.SpringEquinox},
		{model.PureBrightness, y.PureBrightness},
		{model.GrainRain, y.GrainRain},
		{model.StartOfSummer, y.StartOfSummer},
		{model.GrainBuds, y.GrainBuds},
		{model.GrainInEar, y.GrainInEar},
		{model.SummerSolstice, y.SummerSolstice},
		{model.MinorHeat, y.MinorHeat},
		{model.MajorHeat, y.MajorHeat},
		{model.StartOfAutumn, y.StartOfAutumn},
		{model.EndOfHeat, y.EndOfHeat},
		{model.WhiteDew, y.WhiteDew},
		{model.AutumnEquinox, y.AutumnEquinox},
		{model.ColdDew, y.ColdDew},
		{model.Frost, y.Frost},
		{model.StartOfWinter, y.StartOfWinter},
		{model.MinorSnow, y.MinorSnow},
		{model.MajorSnow, y.MajorSnow},
		{model.WinterSolstice, y.WinterSolstice},
	}
	terms := make([]solarTermTime, 0, len(raw))
	for _, r := range raw {
		t, err := time.Parse(solarTermLayout, r.value)
		if err != nil {
			return nil, fmt.Errorf("%w: solar term %s: %v", model.ErrDataUnavailable, r.name, err)
		}
		terms = append(terms, solarTermTime{name: r.name, time: t})
	}
	return terms, nil
}

const solarTermLayout = "2006-01-02 15:04:05.999999999-07:00"

// findSolarTerm locates the solar term containing t. It returns the term name,
// the minutes elapsed since the preceding "initial" term (jie) and the minutes
// remaining until the next one; both bound the current Ba-zi month.
func findSolarTerm(t time.Time, terms, nextYear, previousYear []solarTermTime) (string, int, int, error) {
	// nextInitial resolves the first initial term after index i, which wraps
	// into the following year for the last terms of December.
	nextInitial := func(i int) time.Time {
		idx := i + 1
		if slices.Contains(InitialTerms, terms[i].name) {
			idx = i + 2
		}
		if idx < len(terms) {
			return terms[idx].time
		}
		return nextYear[idx-len(terms)].time
	}

	for i, term := range terms {
		if t.Before(term.time) || (i+1 < len(terms) && !t.Before(terms[i+1].time)) {
			continue
		}

		passed := int(t.Sub(term.time).Minutes())
		if slices.Contains(MidpointTerms, term.name) {
			passed += int(term.time.Sub(terms[i-1].time).Minutes())
		}
		return term.name, passed, int(nextInitial(i).Sub(t).Minutes()), nil
	}

	// t precedes the first term of the year: it still belongs to the month
	// opened by Major Snow of the previous year.
	majorSnow := previousYear[len(previousYear)-2].time
	return model.MajorSnow, int(t.Sub(majorSnow).Minutes()), int(terms[0].time.Sub(t).Minutes()), nil
}

// loadSolarTerms decodes the embedded table once. The returned map is shared
// and must be treated as read-only.
var loadSolarTerms = sync.OnceValues(func() (map[int][]solarTermTime, error) {
	var raw map[string]model.CombinedData
	if err := json.Unmarshal(data.SolarTerm, &raw); err != nil {
		return nil, fmt.Errorf("%w: solar-term.json: %v", model.ErrDataUnavailable, err)
	}
	years := make(map[int][]solarTermTime, len(raw))
	for key, entry := range raw {
		year, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("%w: solar-term.json: year %q", model.ErrDataUnavailable, key)
		}
		terms, err := parseSolarTerms(entry.Data)
		if err != nil {
			return nil, err
		}
		years[year] = terms
	}
	return years, nil
})
