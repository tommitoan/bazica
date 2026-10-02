package analysis

// Star rule tables. Stems are zero-based indices (Jia = 0) and branches are
// Rat-first indices (Rat = 0), the same indices the evaluator uses. Each table is
// transcribed from the classical rule tables and checked against the stars the
// reference page prints; see the tests.

// thienAtTable lists the Heavenly Noble branches by stem.
var thienAtTable = [10][]int{
	{1, 7},
	{0, 8},
	{11, 9},
	{11, 9},
	{1, 7},
	{0, 8},
	{2, 6},
	{2, 6},
	{3, 5},
	{3, 5},
}

// quocAnTable lists the National Seal branch by stem.
var quocAnTable = [10][]int{
	{10},
	{11},
	{1},
	{2},
	{1},
	{2},
	{4},
	{5},
	{7},
	{8},
}

// vanXuongTable lists the Literary Star branch by stem.
var vanXuongTable = [10][]int{
	{5},
	{6},
	{8},
	{9},
	{8},
	{9},
	{11},
	{0},
	{2},
	{3},
}

// phucTinhTable lists the Fortune Star branches by stem.
var phucTinhTable = [10][]int{
	{2, 0},
	{3, 1},
	{0, 2},
	{11},
	{8},
	{7},
	{6},
	{5},
	{4},
	{3, 1},
}

// thaiCucTable lists the Supreme Noble branches by stem.
var thaiCucTable = [10][]int{
	{0, 6},
	{0, 6},
	{3, 9},
	{3, 9},
	{4, 10, 1, 7},
	{4, 10, 1, 7},
	{2, 11},
	{2, 11},
	{5, 8},
	{5, 8},
}

// locTable lists the Prosperity (Loc) branch by stem.
var locTable = [10][]int{
	{2},
	{3},
	{5},
	{6},
	{5},
	{6},
	{8},
	{9},
	{11},
	{0},
}

// kimDuTable lists the Golden Carriage branch by stem.
var kimDuTable = [10][]int{
	{4},
	{5},
	{7},
	{8},
	{7},
	{8},
	{10},
	{11},
	{1},
	{2},
}

// hocDuongTable lists the Study Hall branch by stem.
var hocDuongTable = [10][]int{
	{11},
	{6},
	{2},
	{9},
	{2},
	{9},
	{5},
	{0},
	{8},
	{3},
}

// hongLoanTable is the Red Phoenix branch by year branch.
var hongLoanTable = [12]int{3, 2, 1, 0, 11, 10, 9, 8, 7, 6, 5, 4}

// branchGroups are the four groups of three branches that harmonise (san he).
var branchGroups = [4][3]int{{8, 0, 4}, {11, 3, 7}, {2, 6, 10}, {5, 9, 1}}

// groupStarTargets gives, for each star, its branch in each of the four branch groups.
var (
	vongThanTargets  = [4]int{11, 2, 5, 8}
	kiepSatTargets   = [4]int{5, 8, 11, 2}
	dichMaTargets    = [4]int{2, 5, 8, 11}
	tuongTinhTargets = [4]int{0, 3, 6, 9}
	hoaCaiTargets    = [4]int{4, 7, 10, 1}
	daoHoaTargets    = [4]int{9, 0, 3, 6}
	taiSatTargets    = [4]int{6, 9, 0, 3}
)

// nguyetDucStems is the Monthly Virtue stem for each branch group of the month branch.
var nguyetDucStems = [4]int{8, 0, 2, 6}

// thienDucTable gives the Heavenly Virtue stem or branch by month branch.
var thienDucTable = [12]struct {
	isStem bool
	value  int
}{
	{isStem: false, value: 5},
	{isStem: true, value: 6},
	{isStem: true, value: 3},
	{isStem: false, value: 8},
	{isStem: true, value: 8},
	{isStem: true, value: 7},
	{isStem: false, value: 11},
	{isStem: true, value: 0},
	{isStem: true, value: 9},
	{isStem: false, value: 2},
	{isStem: true, value: 2},
	{isStem: true, value: 1},
}

// thapAcPairs are the day pillars with Ten Evil Great Defeat.
var thapAcPairs = [][2]int{{0, 4}, {1, 5}, {2, 8}, {3, 11}, {4, 10}, {5, 1}, {6, 4}, {7, 5}, {8, 8}, {9, 11}}

// khoiCuongPairs are the pillars with Kui Gang.
var khoiCuongPairs = [][2]int{{4, 10}, {6, 4}, {6, 10}, {8, 4}}

// tienThanPairs are the pillars with Advancing Spirit.
var tienThanPairs = [][2]int{{0, 0}, {0, 6}, {5, 3}, {5, 9}}

// kimThanPairs are the day or hour pillars with Golden Spirit.
var kimThanPairs = [][2]int{{1, 1}, {5, 5}, {9, 9}}

// amDuongSaiThacPairs are the day pillars with Yin Yang Error.
var amDuongSaiThacPairs = [][2]int{{2, 0}, {2, 6}, {3, 1}, {3, 7}, {4, 2}, {4, 8}, {7, 3}, {7, 9}, {8, 4}, {8, 10}, {9, 5}, {9, 11}}

// amDuongSatPairs are the day pillars with Yin Yang Sha.
var amDuongSatPairs = [][2]int{{2, 0}, {4, 6}}

// nhatQuyPairs are the day pillars with Day Noble.
var nhatQuyPairs = [][2]int{{3, 9}, {3, 11}, {9, 3}, {9, 5}}

// lucTuPairs are the day pillars with Six Elegances.
var lucTuPairs = [][2]int{{2, 6}, {3, 7}, {4, 0}, {4, 6}, {5, 1}, {5, 7}}
