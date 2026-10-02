package bazica

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)

// referenceChart returns the chart of the 1995-06-08 11:01 male reference birth as decoded JSON.
func referenceChart(t *testing.T) map[string]any {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	chart, err := GetBaziChart(time.Date(1995, 6, 8, 11, 1, 0, 0, loc), loc, 1)
	if err != nil {
		t.Fatalf("GetBaziChart() error = %v", err)
	}
	return roundTrip(t, chart)
}

func roundTrip(t *testing.T, v any) map[string]any {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func loadFixture(t *testing.T, name string) map[string]any {
	t.Helper()
	data, err := os.ReadFile("testdata/analysis/" + name)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

// stripKeys removes every occurrence of the named keys, at any depth.
func stripKeys(node any, keys ...string) {
	switch n := node.(type) {
	case map[string]any:
		for _, k := range keys {
			delete(n, k)
		}
		for _, v := range n {
			stripKeys(v, keys...)
		}
	case []any:
		for _, v := range n {
			stripKeys(v, keys...)
		}
	}
}

// The fixture is the contract shared with bazica-web.
func TestAnalysisMatchesTheSharedFixture(t *testing.T) {
	got := referenceChart(t)
	want := loadFixture(t, "chart_1995_expected_v140.json")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("chart differs from fixture near:\n%s", firstDifference(mustJSON(t, got), mustJSON(t, want)))
	}
}

// firstDifference returns a short excerpt around the first differing byte.
func firstDifference(a, b []byte) string {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			from, to := max(0, i-120), min(len(a), i+120)
			return string(a[from:to])
		}
	}
	return "lengths differ"
}

// Existing consumers read the v1.3.0 fields; they must be byte-for-byte equal
// once the added analysis objects are removed.
func TestLegacyFieldsAreUnchangedByAnalysis(t *testing.T) {
	got := referenceChart(t)
	stripKeys(got, "analysis")
	want := loadFixture(t, "chart_1995_v130.json")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("v1.3.0 fields changed:\n%s", firstDifference(mustJSON(t, got), mustJSON(t, want)))
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestAnalysisSerialisesNullsAndLists(t *testing.T) {
	chart := referenceChart(t)
	fp := chart["four_pillars"].(map[string]any)
	day := fp["day_pillar"].(map[string]any)["analysis"].(map[string]any)
	if v, ok := day["ten_god"]; !ok || v != nil {
		t.Errorf("day ten_god = %v (present %v), want null", v, ok)
	}
	year := fp["year_pillar"].(map[string]any)["analysis"].(map[string]any)
	if stars, ok := year["stars"].([]any); !ok || len(stars) == 0 {
		t.Errorf("year stars = %#v, want the stars of the 1995 chart", year["stars"])
	}
	analysis := chart["analysis"].(map[string]any)
	for _, k := range []string{"day_master_strength", "useful_god"} {
		if v, ok := analysis[k]; !ok || v != nil {
			t.Errorf("%s = %v (present %v), want null", k, v, ok)
		}
	}
}

func TestAnalysisIsAttachedForPillarsAcrossTheCalendar(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	for _, d := range []time.Time{
		time.Date(1900, 1, 1, 0, 30, 0, 0, loc),
		time.Date(1984, 7, 4, 23, 30, 0, 0, loc),
		time.Date(2023, 2, 4, 9, 55, 0, 0, loc),
		time.Date(2099, 12, 31, 22, 59, 0, 0, loc),
	} {
		for gender := 0; gender <= 1; gender++ {
			chart, err := GetBaziChart(d, loc, gender)
			if err != nil {
				t.Fatalf("%v: %v", d, err)
			}
			fp := chart.FourPillar
			if fp.YearPillar.Analysis == nil || fp.MonthPillar.Analysis == nil || fp.DayPillar.Analysis == nil || fp.HourPillar.Analysis == nil || chart.Analysis == nil {
				t.Errorf("%v: analysis missing", d)
			}
		}
	}
}
