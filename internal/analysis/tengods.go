package analysis

import "github.com/tommitoan/bazica/model"

// tenGod classifies stem against the Day Master. The element relation picks the
// pair of gods (same element, generated, controlled, controlling, generating)
// and the polarity picks the first name for equal polarity or the second otherwise.
func tenGod(dayStem, stem int) model.LocalizedTerm {
	relation := (stemElement(stem) - stemElement(dayStem) + 5) % 5
	polarity := 0
	if isYang(dayStem) != isYang(stem) {
		polarity = 1
	}
	return tenGodTerms[relation][polarity]
}
