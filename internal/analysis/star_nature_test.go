package analysis

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/tommitoan/bazica/model"
)

var validNatures = map[string]bool{
	model.NatureAuspicious:   true,
	model.NatureInauspicious: true,
	model.NatureMixed:        true,
}

// Every registered star has exactly one valid nature, and no nature is left over for a star that does not exist.
func TestEveryStarHasAValidNature(t *testing.T) {
	for code, term := range starTerms {
		nature, ok := starNature[code]
		if !ok {
			t.Errorf("star %s has no entry in starNature", code)
			continue
		}
		if !validNatures[nature] {
			t.Errorf("star %s has the unknown nature %q", code, nature)
		}
		if term.Nature != nature {
			t.Errorf("star %s: term carries %q, table says %q", code, term.Nature, nature)
		}
	}
	for code := range starNature {
		if _, ok := starTerms[code]; !ok {
			t.Errorf("starNature lists %s, which is not a registered star", code)
		}
	}
}

// The classification must not drift into the other terms: Ten Gods, stages, elements, polarities and Nayin
// never carry a nature.
func TestOnlyStarsCarryANature(t *testing.T) {
	for _, term := range allTerms() {
		if term.Nature != "" {
			t.Errorf("term %s carries the nature %q", term.Code, term.Nature)
		}
		data, err := json.Marshal(term)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "nature") {
			t.Errorf("term %s serialises a nature: %s", term.Code, data)
		}
	}
}

// The three levels each have stars, so a typo cannot silently collapse the table into one level.
func TestAllThreeNaturesAreUsed(t *testing.T) {
	seen := map[string]int{}
	for _, nature := range starNature {
		seen[nature]++
	}
	for nature := range validNatures {
		if seen[nature] == 0 {
			t.Errorf("no star is %s", nature)
		}
	}
}

// The star JSON gains one key and loses none.
func TestStarJSONShape(t *testing.T) {
	data, err := json.Marshal(starTerms["quoc_an"])
	if err != nil {
		t.Fatal(err)
	}
	want := `{"code":"star.quoc_an","en":"National Seal","vi":"Quốc Ấn","nature":"auspicious"}`
	if string(data) != want {
		t.Errorf("star JSON = %s, want %s", data, want)
	}
}
