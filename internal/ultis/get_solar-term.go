package ultis

import (
	"encoding/json"
	"fmt"
	"github.com/tommitoan/bazica/model"
	"os"
	"slices"
	"time"
)

var InititalTerms = []string{
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

func GetSolarTerm(path string, dateTime time.Time) (string, int, int, error) {
	year := dateTime.Year()
	current, next, previous, err := GetSolarTermsByYear(fmt.Sprint(year), fmt.Sprint(year+1), fmt.Sprint(year-1), path)
	if err != nil {
		return "", 0, 0, err
	}
	return findSolarTerm(dateTime, current, next, previous)
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
func findSolarTerm(t time.Time, currentYear, nextYear, previousYear model.SolarTermYear) (string, int, int, error) {
	terms, err := parseSolarTerms(currentYear)
	if err != nil {
		return "", 0, 0, err
	}

	// nextInitial resolves the first initial term after index i, which wraps
	// into the following year for the last terms of December.
	nextInitial := func(i int) (time.Time, error) {
		idx := i + 1
		if slices.Contains(InititalTerms, terms[i].name) {
			idx = i + 2
		}
		if idx < len(terms) {
			return terms[idx].time, nil
		}
		following, err := parseSolarTerms(nextYear)
		if err != nil {
			return time.Time{}, err
		}
		return following[idx-len(terms)].time, nil
	}

	for i, term := range terms {
		if t.Before(term.time) || (i+1 < len(terms) && !t.Before(terms[i+1].time)) {
			continue
		}

		passed := int(t.Sub(term.time).Minutes())
		if slices.Contains(MidpointTerms, term.name) {
			passed += int(term.time.Sub(terms[i-1].time).Minutes())
		}
		boundary, err := nextInitial(i)
		if err != nil {
			return "", 0, 0, err
		}
		return term.name, passed, int(boundary.Sub(t).Minutes()), nil
	}

	// t precedes the first term of the year: it still belongs to the month
	// opened by Major Snow of the previous year.
	previous, err := parseSolarTerms(previousYear)
	if err != nil {
		return "", 0, 0, err
	}
	majorSnow := previous[len(previous)-2].time
	passed := int(t.Sub(majorSnow).Minutes())
	remaining := int(terms[0].time.Sub(t).Minutes())
	return model.MajorSnow, passed, remaining, nil
}

var PrefixPath string

// GetSolarTermsByYear returns the solar term data of the requested year and its
// neighbours. A year without data yields model.ErrDateOutOfRange.
func GetSolarTermsByYear(year, nextYear, previousYear string, path ...string) (model.SolarTermYear, model.SolarTermYear, model.SolarTermYear, error) {
	var prefix string
	if len(path) != 0 {
		prefix = path[0]
	}
	data, err := getSolarTermData(prefix)
	if err != nil {
		return model.SolarTermYear{}, model.SolarTermYear{}, model.SolarTermYear{}, err
	}
	var out [3]model.SolarTermYear
	for i, y := range []string{year, nextYear, previousYear} {
		entry, ok := data[y]
		if !ok {
			return model.SolarTermYear{}, model.SolarTermYear{}, model.SolarTermYear{}, fmt.Errorf("%w: no solar term data for %s", model.ErrDateOutOfRange, y)
		}
		out[i] = entry.Data
	}
	return out[0], out[1], out[2], nil
}

func getSolarTermData(path string) (map[string]model.CombinedData, error) {
	prefix := PrefixPath
	if path != "" {
		prefix = path
	}

	data, err := os.ReadFile(prefix + "data/solar-term.json")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrDataUnavailable, err)
	}

	if err := json.Unmarshal(data, &model.TempCombinedData); err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrDataUnavailable, err)
	}
	if model.TempCombinedData == nil {
		return nil, fmt.Errorf("%w: solar-term.json is empty", model.ErrDataUnavailable)
	}
	return model.TempCombinedData, nil
}
