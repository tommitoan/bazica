package analysis

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/tommitoan/bazica/v2/model"
)

type registryFile struct {
	Terms []model.LocalizedTerm `json:"terms"`
}

func allTerms() []model.LocalizedTerm {
	var terms []model.LocalizedTerm
	for _, pair := range tenGodTerms {
		terms = append(terms, pair[0], pair[1])
	}
	terms = append(terms, stageTerms[:]...)
	terms = append(terms, elementTerms[:]...)
	terms = append(terms, polarityTerms[:]...)
	for _, t := range nayinTerms {
		terms = append(terms, t)
	}
	return terms
}

// The shared registry fixture also lists stars, which are not part of the
// tables yet; every other term must match it exactly.
func TestTermTablesMatchTheSharedRegistry(t *testing.T) {
	data, err := os.ReadFile("../../testdata/analysis/term_registry.json")
	if err != nil {
		t.Skipf("registry fixture not found: %v", err)
	}
	var file registryFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	want := map[string]model.LocalizedTerm{}
	for _, term := range file.Terms {
		if !strings.HasPrefix(term.Code, "star.") {
			want[term.Code] = term
		}
	}
	got := map[string]model.LocalizedTerm{}
	for _, term := range allTerms() {
		if _, dup := got[term.Code]; dup {
			t.Errorf("duplicate code %s", term.Code)
		}
		got[term.Code] = term
	}
	if len(got) != len(want) {
		t.Errorf("table has %d terms, registry has %d", len(got), len(want))
	}
	for code, w := range want {
		if got[code] != w {
			t.Errorf("%s: got %+v, want %+v", code, got[code], w)
		}
	}
}

func TestEveryTermHasCodeAndBothLabels(t *testing.T) {
	for _, term := range allTerms() {
		if term.Code == "" || term.EN == "" || term.VI == "" {
			t.Errorf("incomplete term %+v", term)
		}
	}
}
