package analysis

// hiddenStems lists the stems hidden in each branch, main qi first, as stem
// values. It is indexed by branchIndex (Rat first). No weights are applied.
var hiddenStems = [12][]int{
	{10}, {6, 10, 8}, {1, 3, 5}, {2}, {5, 2, 10}, {3, 5, 7},
	{4, 6}, {6, 4, 2}, {7, 9, 5}, {8}, {5, 8, 4}, {9, 1},
}

func hiddenStemsOf(branch int) []int { return hiddenStems[branchIndex(branch)] }
