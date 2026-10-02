package analysis

// voidBranches returns the two branch values missing from the ten-pillar cycle
// of the Day pillar: the cycle starts at the Jia branch found by walking back
// by the day stem index, and the last two branches of the twelve are not paired.
func voidBranches(dayStem, dayBranch int) [2]int {
	first := branchIndex(dayBranch) - stemIndex(dayStem)
	return [2]int{branchValue(first + 10), branchValue(first + 11)}
}
