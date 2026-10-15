package analysis

import (
	"strings"
	"testing"

	"github.com/tommitoan/bazica/v2/internal/utils"
)

func TestSolarTermTermsCoverAllCalendarTerms(t *testing.T) {
	seen := make(map[string]bool)
	for _, name := range append(append([]string(nil), utils.InitialTerms...), utils.MidpointTerms...) {
		term, ok := solarTermTerms[name]
		if !ok {
			t.Errorf("missing term %s", name)
			continue
		}
		if term.Code != "solar_term."+name || term.EN == "" || term.VI == "" {
			t.Errorf("invalid label for %s: %+v", name, term)
		}
		if strings.TrimSpace(term.EN) != term.EN || strings.TrimSpace(term.VI) != term.VI {
			t.Errorf("term %s has surrounding whitespace", name)
		}
		if seen[name] {
			t.Errorf("duplicate term %s", name)
		}
		seen[name] = true
	}
	if len(seen) != 24 || len(solarTermTerms) != 24 {
		t.Errorf("registry = %d terms, calendar = %d, want 24 each", len(solarTermTerms), len(seen))
	}
}
