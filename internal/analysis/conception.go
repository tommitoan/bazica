package analysis

// thaiNguyen returns the conception pillar: month stem + 1 and month branch + 3.
func thaiNguyen(monthStem, monthBranch int) (stem, branch int) {
	return stemValue(stemIndex(monthStem) + 1), branchValue(branchIndex(monthBranch) + 3)
}

// thaiTuc returns the gestation pillar: the five-harmony partner of the day
// stem and the six-harmony partner of the day branch.
func thaiTuc(dayStem, dayBranch int) (stem, branch int) {
	return stemValue(stemIndex(dayStem) + 5), branchValue(1 - branchIndex(dayBranch))
}

// lifePalace returns the life palace pillar. The branch is (5 - month - hour)
// counted Rat first. The stem follows the Five Tigers stem of the Tiger branch
// for the year stem, advanced by the branch offset from Tiger; Rat and Ox are
// counted before Tiger (offsets -2 and -1), not after Pig as in the classical
// month order. That Rat and Ox rule is the convention of the reference page.
func lifePalace(yearStem, monthBranch, hourBranch int) (stem, branch int) {
	b := ((5-branchIndex(monthBranch)-branchIndex(hourBranch))%12 + 12) % 12
	tigerStem := (stemIndex(yearStem)%5)*2 + 2
	return stemValue(tigerStem + b - 2), branchValue(b)
}
