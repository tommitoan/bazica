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

	inAll          = inYear | inMonth | inDay | inHour
	inMonthHour    = inMonth | inHour
	inMonthDayHour = inMonth | inDay | inHour
	inDayHour      = inDay | inHour
	inNotDay       = inYear | inMonth | inHour
)

func (s pillarSet) has(pillar int) bool { return s&(1<<pillar) != 0 }

// chartFacts is what the star rules read: zero-based stem indices, Rat-first
// branch indices and the gender.
type chartFacts struct {
	stem   [4]int
	branch [4]int
	male   bool
	// general is the branch of the Moon General at the birth date; it is only
	// meaningful when hasGeneral is set.
	general    int
	hasGeneral bool
}

// starRule is one row of the star table: the star it awards, the pillars it
// may land in, and the test for one pillar.
type starRule struct {
	code  string
	where pillarSet
	match func(f *chartFacts, pillar int) bool
}

// newChartFacts reads the pillars. general is the branch of the Moon General,
// or a negative number when the birth date is unknown.
func newChartFacts(pillars [4]natalPillar, male bool, general int) *chartFacts {
	f := &chartFacts{male: male, general: general, hasGeneral: general >= 0}
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

// daiHao awards Great Depletion: the branch clashing with the year branch, one
// step on for a Yang-year man or a Yin-year woman and one step back otherwise
// (the same direction rule as Giao and Cau).
func daiHao(code string) starRule {
	return starRule{code, inAll, func(f *chartFacts, p int) bool {
		step := -1
		if (f.stem[yearPillar]%2 == 0) == f.male {
			step = 1
		}
		return f.branch[p] == mod12(opposite(f.branch[yearPillar])+step)
	}}
}

// byMonthGroupStem awards a star to a pillar whose stem is listed for the
// branch group of the month branch.
func byMonthGroupStem(code string, stems [4][]int) starRule {
	return starRule{code, inAll, func(f *chartFacts, p int) bool {
		return contains(stems[groupOf(f.branch[monthPillar])], f.stem[p])
	}}
}

// seasonIndex returns which of the four season groups holds a branch.
func seasonIndex(groups [4][3]int, branch int) int {
	for i, g := range groups {
		if contains(g[:], branch) {
			return i
		}
	}
	return 0
}

// bySeasonGroup awards a star on the month, day and hour pillars whose branch
// is the target of the season group of the year branch.
func bySeasonGroup(code string, targets [4]int) starRule {
	return starRule{code, inMonthDayHour, func(f *chartFacts, p int) bool {
		return f.branch[p] == targets[seasonIndex(seasonGroups, f.branch[yearPillar])]
	}}
}

// byPairedDayBranch awards a star to a pillar other than the day pillar whose
// branch is one of two branches while the day branch is the other one.
func byPairedDayBranch(code string, pair [2]int) starRule {
	return starRule{code, inNotDay, func(f *chartFacts, p int) bool {
		day := f.branch[dayPillar]
		return contains(pair[:], f.branch[p]) && contains(pair[:], day) && day != f.branch[p]
	}}
}

// bySeasonDayPillar awards a star to a day pillar listed for the season of the month branch.
func bySeasonDayPillar(code string, pairs [4][][2]int) starRule {
	return starRule{code, onlyPillar(dayPillar), func(f *chartFacts, p int) bool {
		for _, pair := range pairs[seasonIndex(seasonOfMonthBranch, f.branch[monthPillar])] {
			if pair[0] == f.stem[p] && pair[1] == f.branch[p] {
				return true
			}
		}
		return false
	}}
}

// tamKy awards a Three Wonders star to three consecutive pillars, counted from
// the year pillar, whose stems are the triple in order.
func tamKy(code string, triple [3]int) starRule {
	return starRule{code, inAll, func(f *chartFacts, p int) bool {
		for start := 0; start <= 1; start++ {
			if p < start || p > start+2 {
				continue
			}
			if f.stem[start] == triple[0] && f.stem[start+1] == triple[1] && f.stem[start+2] == triple[2] {
				return true
			}
		}
		return false
	}}
}

// moonGeneral awards the Moon General to every pillar whose branch is the
// general of the birth date. Without a birth date nothing is awarded.
func moonGeneral(code string) starRule {
	return starRule{code, inAll, func(f *chartFacts, p int) bool {
		return f.hasGeneral && f.branch[p] == f.general
	}}
}
