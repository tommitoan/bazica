package analysis

import "github.com/tommitoan/bazica/v2/model"

// Element indices follow the generating cycle: Wood, Fire, Earth, Metal, Water.
var (
	stemElements = [10]int{0, 0, 1, 1, 2, 2, 3, 3, 4, 4}
	// branchElements is indexed by branchIndex (Rat first).
	branchElements = [12]int{4, 2, 0, 0, 2, 1, 1, 2, 3, 3, 2, 4}
)

func stemElement(stem int) int { return stemElements[stemIndex(stem)] }

func branchElement(branch int) int { return branchElements[branchIndex(branch)] }

// isYang reports the polarity of a stem; odd stem values (Jia, Bing, ...) are Yang.
func isYang(stem int) bool { return stem%2 == 1 }

func polarityTerm(stem int) model.LocalizedTerm {
	if isYang(stem) {
		return polarityTerms[0]
	}
	return polarityTerms[1]
}

func addElement(c *model.ElementCount, element int) {
	switch element {
	case 0:
		c.Wood++
	case 1:
		c.Fire++
	case 2:
		c.Earth++
	case 3:
		c.Metal++
	case 4:
		c.Water++
	}
}
