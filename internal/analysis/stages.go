package analysis

import (
	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

// stage returns the life stage of any stem at a branch.
func stage(stem, branch int) model.LocalizedTerm {
	return stageTerms[utils.GetLifeCycleValueOfStem(stem, branch)-1]
}
