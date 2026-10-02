package analysis

import "github.com/tommitoan/bazica/model"

// Pillar positions in a chart.
const (
	yearPillar = iota
	monthPillar
	dayPillar
	hourPillar
)

// pillarSet is a set of pillar positions.
type pillarSet uint8

const (
	inYear pillarSet = 1 << iota
	inMonth
	inDay
	inHour

	inAll       = inYear | inMonth | inDay | inHour
	inMonthHour = inMonth | inHour
	inDayHour   = inDay | inHour
	inNotDay    = inYear | inMonth | inHour
)

func (s pillarSet) has(pillar int) bool { return s&(1<<pillar) != 0 }

// chartFacts is what the star rules read: zero-based stem indices, Rat-first
// branch indices and the gender.
type chartFacts struct {
	stem   [4]int
	branch [4]int
	male   bool
}

// starRule is one row of the star table: the star it awards, the pillars it
// may land in, and the test for one pillar.
type starRule struct {
	code  string
	where pillarSet
	match func(f *chartFacts, pillar int) bool
}

func newChartFacts(pillars [4]natalPillar, male bool) *chartFacts {
	f := &chartFacts{male: male}
	for i, p := range pillars {
		f.stem[i] = stemIndex(p.stem.Value)
		f.branch[i] = branchIndex(p.branch.Value)
	}
	return f
}

// evaluateStars returns the stars of each pillar in registry order. A pillar
// without stars gets an empty, non-nil list.
func evaluateStars(f *chartFacts) [4][]model.LocalizedTerm {
	var out [4][]model.LocalizedTerm
	for pillar := range out {
		out[pillar] = []model.LocalizedTerm{}
		for _, rule := range starRules {
			if rule.where.has(pillar) && rule.match(f, pillar) {
				out[pillar] = append(out[pillar], starTerms[rule.code])
			}
		}
	}
	return out
}

func contains(list []int, v int) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func mod12(n int) int { return ((n % 12) + 12) % 12 }

// opposite is the clashing branch, six places away.
func opposite(branch int) int { return mod12(branch + 6) }

// harmony is the six-harmony partner of a branch (Rat with Ox, Tiger with Pig, ...).
func harmony(branch int) int { return mod12(1 - branch) }

// groupOf returns which of the four branch groups of three contains branch.
func groupOf(branch int) int {
	for i, group := range branchGroups {
		if contains(group[:], branch) {
			return i
		}
	}
	return 0
}

// voidBranchIndices returns the Rat-first indices of the two void branches of a day pillar.
func voidBranchIndices(dayStem, dayBranch int) [2]int {
	first := mod12(dayBranch - dayStem)
	return [2]int{mod12(first + 10), mod12(first + 11)}
}

// Bases of a rule.
var (
	dayOrYear = []int{dayPillar, yearPillar}
	dayOnly   = []int{dayPillar}
	yearOnly  = []int{yearPillar}
	yearOrDay = []int{yearPillar, dayPillar}
)

// byStem awards a star to a pillar whose branch is listed for the stem of any base pillar.
func byStem(code string, bases []int, table [10][]int, where pillarSet) starRule {
	return starRule{code, where, func(f *chartFacts, p int) bool {
		for _, base := range bases {
			if contains(table[f.stem[base]], f.branch[p]) {
				return true
			}
		}
		return false
	}}
}

// byGroup awards a star to a pillar whose branch is the target of the branch
// group of any base pillar's branch.
func byGroup(code string, bases []int, targets [4]int, where pillarSet) starRule {
	return starRule{code, where, func(f *chartFacts, p int) bool {
		for _, base := range bases {
			if f.branch[p] == targets[groupOf(f.branch[base])] {
				return true
			}
		}
		return false
	}}
}

// byBranch awards a star to a pillar whose branch is listed for the branch of the base pillar.
func byBranch(code string, base int, table [12][]int, where pillarSet) starRule {
	return starRule{code, where, func(f *chartFacts, p int) bool {
		return contains(table[f.branch[base]], f.branch[p])
	}}
}

// byPair awards a star to a pillar whose stem and branch pair is listed.
func byPair(code string, pairs [][2]int, where pillarSet) starRule {
	return starRule{code, where, func(f *chartFacts, p int) bool {
		for _, pair := range pairs {
			if pair[0] == f.stem[p] && pair[1] == f.branch[p] {
				return true
			}
		}
		return false
	}}
}

// onlyPillar is the set holding one pillar position.
func onlyPillar(pillar int) pillarSet { return 1 << pillar }

// deriveStemTable builds a stem table from another by mapping each row.
func deriveStemTable(src [10][]int, fn func(stem int, branches []int) []int) [10][]int {
	var out [10][]int
	for stem, branches := range src {
		out[stem] = fn(stem, branches)
	}
	return out
}

// deriveBranchTable builds a table indexed by branch from a function.
func deriveBranchTable(fn func(branch int) []int) [12][]int {
	var out [12][]int
	for branch := range out {
		out[branch] = fn(branch)
	}
	return out
}

// giaoCau awards Giao or Cau: three branches ahead of or behind the year
// branch. For a Yang year man or a Yin year woman Cau is ahead and Giao behind;
// otherwise they swap.
func giaoCau(code string, isCau bool) starRule {
	return starRule{code, inAll, func(f *chartFacts, p int) bool {
		year := f.branch[yearPillar]
		ahead, behind := mod12(year+3), mod12(year-3)
		forward := (f.stem[yearPillar]%2 == 0) == f.male
		if forward == isCau {
			return f.branch[p] == ahead
		}
		return f.branch[p] == behind
	}}
}

// nguyetDuc awards Monthly Virtue to a pillar whose stem is the stem of the month branch's group.
func nguyetDuc(code string) starRule {
	return starRule{code, inAll, func(f *chartFacts, p int) bool {
		return f.stem[p] == nguyetDucStems[groupOf(f.branch[monthPillar])]
	}}
}

// thienDuc awards Heavenly Virtue: a stem or a branch depending on the month branch.
func thienDuc(code string) starRule {
	return starRule{code, inAll, func(f *chartFacts, p int) bool {
		entry := thienDucTable[f.branch[monthPillar]]
		if entry.isStem {
			return f.stem[p] == entry.value
		}
		return f.branch[p] == entry.value
	}}
}

// voidStar marks the pillars whose branch is void in the day pillar's cycle.
func voidStar(code string) starRule {
	return starRule{code, inNotDay, func(f *chartFacts, p int) bool {
		void := voidBranchIndices(f.stem[dayPillar], f.branch[dayPillar])
		return f.branch[p] == void[0] || f.branch[p] == void[1]
	}}
}
