package analysis

// Stems and branches enter this package with the public numbering of the model
// (stem 1 = Yang Wood to 10 = Yin Water, branch 1 = Tiger to 12 = Ox). The
// arithmetic below is easier with Rat-first zero-based indices, so every
// crossing goes through these four helpers.

func stemIndex(value int) int { return value - 1 }

func stemValue(index int) int { return ((index%10)+10)%10 + 1 }

func branchIndex(value int) int { return (value + 1) % 12 }

func branchValue(index int) int { return ((index-2)%12+12)%12 + 1 }

func validStem(value int) bool { return value >= 1 && value <= 10 }

func validBranch(value int) bool { return value >= 1 && value <= 12 }
