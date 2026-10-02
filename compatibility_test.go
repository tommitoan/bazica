package bazica

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

type legacyCase struct {
	Year, Month, Day, Hour, Minute, Gender int
	Chart                                  map[string]any
}

func loadLegacyCases(t *testing.T) []legacyCase {
	t.Helper()
	data, err := os.ReadFile("testdata/analysis/legacy_v130_charts.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var file struct {
		Charts []struct {
			Year   int            `json:"year"`
			Month  int            `json:"month"`
			Day    int            `json:"day"`
			Hour   int            `json:"hour"`
			Minute int            `json:"minute"`
			Gender int            `json:"gender"`
			Chart  map[string]any `json:"chart"`
		} `json:"charts"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	cases := make([]legacyCase, 0, len(file.Charts))
	for _, c := range file.Charts {
		cases = append(cases, legacyCase{c.Year, c.Month, c.Day, c.Hour, c.Minute, c.Gender, c.Chart})
	}
	return cases
}

func chartAt(t *testing.T, c legacyCase) map[string]any {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	chart, err := GetBaziChart(time.Date(c.Year, time.Month(c.Month), c.Day, c.Hour, c.Minute, 0, 0, loc), loc, c.Gender)
	if err != nil {
		t.Fatalf("%d-%02d-%02d: %v", c.Year, c.Month, c.Day, err)
	}
	return roundTrip(t, chart)
}

// Charts computed by the published v1.3.0 module and the current library must
// differ only by the added analysis objects.
func TestChartsDifferFromV130OnlyByAddedFields(t *testing.T) {
	cases := loadLegacyCases(t)
	if len(cases) < 8 {
		t.Fatalf("only %d pinned charts", len(cases))
	}
	for _, c := range cases {
		got := chartAt(t, c)
		stripKeys(got, "analysis")
		if !reflect.DeepEqual(got, c.Chart) {
			t.Errorf("%d-%02d-%02d %02d:%02d gender %d differs from v1.3.0:\n%s", c.Year, c.Month, c.Day, c.Hour, c.Minute, c.Gender,
				firstDifference(mustJSON(t, got), mustJSON(t, c.Chart)))
		}
	}
}

// Every added key must be named analysis, so an old client that ignores unknown keys is unaffected.
func TestEveryAddedKeyIsNamedAnalysis(t *testing.T) {
	for _, c := range loadLegacyCases(t) {
		if extra := addedKeys(chartAt(t, c), c.Chart, ""); len(extra) != 0 {
			t.Errorf("%d-%02d-%02d: keys added outside analysis: %v", c.Year, c.Month, c.Day, extra)
		}
	}
}

// addedKeys lists the paths present in got but not in want, ignoring analysis objects.
func addedKeys(got, want any, path string) []string {
	var added []string
	switch g := got.(type) {
	case map[string]any:
		w, _ := want.(map[string]any)
		for k, v := range g {
			if k == "analysis" {
				continue
			}
			wv, ok := w[k]
			if !ok {
				added = append(added, path+"/"+k)
				continue
			}
			added = append(added, addedKeys(v, wv, path+"/"+k)...)
		}
	case []any:
		w, _ := want.([]any)
		for i, v := range g {
			if i < len(w) {
				added = append(added, addedKeys(v, w[i], path+"/[]")...)
			}
		}
	}
	return added
}

// webView is the part of the response bazica-web and its script read: the four
// pillars and the luck pillars. Decoding the new output into it must give the
// same values as decoding v1.3.0.
type webView struct {
	FourPillars map[string]struct {
		HeavenlyStem  map[string]any `json:"heavenly_stem"`
		EarthlyBranch map[string]any `json:"earthly_branch"`
		GanZhi        map[string]any `json:"gan_zhi"`
		LifeCycle     string         `json:"life_cycle"`
	} `json:"four_pillars"`
	LuckPillars struct {
		LuckPillars []struct {
			Number        int            `json:"number"`
			HeavenlyStem  map[string]any `json:"heavenly_stem"`
			EarthlyBranch map[string]any `json:"earthly_branch"`
			GanZhi        map[string]any `json:"gan_zhi"`
			YearStart     int            `json:"year_start"`
			YearEnd       int            `json:"year_end"`
		} `json:"luck_pillars"`
	} `json:"luck_pillars"`
}

func decodeWebView(t *testing.T, chart map[string]any) webView {
	t.Helper()
	var v webView
	if err := json.Unmarshal(mustJSON(t, chart), &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestWebConsumerReadsTheSameValuesAsV130(t *testing.T) {
	for _, c := range loadLegacyCases(t) {
		if got, want := decodeWebView(t, chartAt(t, c)), decodeWebView(t, c.Chart); !reflect.DeepEqual(got, want) {
			t.Errorf("%d-%02d-%02d: the web view of the chart changed", c.Year, c.Month, c.Day)
		}
	}
}

func TestAddedKeysDetectsAnUnexpectedKey(t *testing.T) {
	want := map[string]any{"four_pillars": map[string]any{"year_pillar": map[string]any{"year": 1995.0}}}
	got := map[string]any{"four_pillars": map[string]any{"year_pillar": map[string]any{"year": 1995.0, "surprise": 1.0, "analysis": map[string]any{}}}}
	extra := addedKeys(got, want, "")
	if len(extra) != 1 || extra[0] != "/four_pillars/year_pillar/surprise" {
		t.Errorf("addedKeys = %v, want only the surprise key", extra)
	}
}

// collectTerms gathers every {code, en, vi} object in a decoded JSON value.
func collectTerms(node any, into map[string]map[string]bool) {
	switch n := node.(type) {
	case map[string]any:
		if code, ok := n["code"].(string); ok && len(n) == 3 {
			label := fmt.Sprintf("%v | %v", n["en"], n["vi"])
			if into[code] == nil {
				into[code] = map[string]bool{}
			}
			into[code][label] = true
			return
		}
		for _, v := range n {
			collectTerms(v, into)
		}
	case []any:
		for _, v := range n {
			collectTerms(v, into)
		}
	}
}

// A term code must carry one label pair wherever it appears: natal pillars,
// luck pillars, auxiliary pillars and the yearly table.
func TestTermCodesKeepOneLabelPairAcrossAllOutput(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	seen := map[string]map[string]bool{}
	for _, c := range loadLegacyCases(t) {
		chart, err := GetBaziChart(time.Date(c.Year, time.Month(c.Month), c.Day, c.Hour, c.Minute, 0, 0, loc), loc, c.Gender)
		if err != nil {
			t.Fatal(err)
		}
		collectTerms(roundTrip(t, chart), seen)
		from := chart.FourPillar.YearPillar.Year - 1
		if from < 1900 {
			from = 1900
		}
		annual, err := GetAnnualPillars(chart, from, 2099-from+1)
		if err != nil {
			annual, err = GetAnnualPillars(chart, chart.FourPillar.YearPillar.Year, 2099-chart.FourPillar.YearPillar.Year+1)
		}
		if err != nil {
			t.Fatal(err)
		}
		collectTerms(roundTrip(t, annual), seen)
	}
	if len(seen) < 40 {
		t.Fatalf("only %d distinct codes seen", len(seen))
	}
	for code, labels := range seen {
		if len(labels) != 1 {
			t.Errorf("%s has %d label pairs: %v", code, len(labels), labels)
		}
	}
}
