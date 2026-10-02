package analysis

// cell is a stem and branch pair taken by value, as used by the clash rule.
type cell struct{ stem, branch int }

// overcomes reports whether stem overcomes other with the same polarity: the
// element of other follows the element of stem by two steps of the generating cycle.
func overcomes(stem, other int) bool {
	return (stemElement(other)-stemElement(stem)+5)%5 == 2 && isYang(stem) == isYang(other)
}

// heavenEarthClash reports whether c's stem overcomes the stem of any other
// cell (same polarity) while their branches are opposite. The relation is
// directional: the overcome cell is not flagged by this pair alone. Earth
// stems take part, so the rule is wider than the four classical pairs.
func heavenEarthClash(c cell, others []cell) bool {
	for _, o := range others {
		if overcomes(c.stem, o.stem) && ((branchIndex(c.branch)-branchIndex(o.branch))%12+12)%12 == 6 {
			return true
		}
	}
	return false
}
