package analysis

// Tables of the stars resolved in the second pass (planning/07). Same indices as star_tables.go.

// hongDiemTable lists the Red Allure branches by stem; the day stem and the year stem both apply.
var hongDiemTable = [10][]int{
	{6, 8},
	{6},
	{2},
	{7},
	{4, 0},
	{4},
	{10},
	{9},
	{0, 5},
	{8},
}

// hocSyTable lists the Scholar branch by day stem.
var hocSyTable = [10][]int{
	{0},
	{11},
	{3},
	{2},
	{6},
	{5},
	{6},
	{5},
	{9},
	{8},
}

// seasonGroups are the year-branch groups of Lonely Spirit and Widowhood: winter (Pig, Rat, Ox), spring, summer and autumn.
var seasonGroups = [4][3]int{{11, 0, 1}, {2, 3, 4}, {5, 6, 7}, {8, 9, 10}}

// seasonOfMonthBranch lists the month branches of spring, summer, autumn and winter, in that order.
var seasonOfMonthBranch = [4][3]int{{2, 3, 4}, {5, 6, 7}, {8, 9, 10}, {11, 0, 1}}

// ducQuyNhanStems lists, by branch group of the month branch, the stems that carry Virtue Noble.
var ducQuyNhanStems = [4][]int{{4, 5, 8, 9}, {0, 1}, {}, {6, 7}}

// tuQuyNhanStems lists, by branch group of the month branch, the stems that carry Elegance Noble.
var tuQuyNhanStems = [4][]int{{}, {3, 8}, {2, 3}, {1}}

// coThanTargets and quaTuTargets give the branch of Lonely Spirit and Widowhood by season group of the year branch.
var (
	coThanTargets = [4]int{2, 1, 8, 11}
	quaTuTargets  = [4]int{10, 5, 4, 7}
)

// coLoanSatPairs are the day pillars with Lone Phoenix Sha.
var coLoanSatPairs = [][2]int{{0, 2}, {1, 5}, {2, 6}, {3, 5}, {4, 6}, {4, 8}, {8, 0}, {9, 5}}

// nhatNhanPairs are the day pillars with Day Blade.
var nhatNhanPairs = [][2]int{{2, 6}, {3, 5}, {4, 6}, {5, 5}, {8, 0}, {9, 11}}

// cuuQuyPhongHaiPairs are the day or hour pillars with Nine Ghosts Harm.
var cuuQuyPhongHaiPairs = [][2]int{{3, 3}, {3, 9}, {4, 0}, {4, 6}, {5, 3}, {5, 9}, {7, 3}, {7, 9}, {8, 0}, {8, 6}}

// thienDiaChuyenSatPairs are the day pillars with Heaven Earth Turning Sha, by season of the month branch.
var thienDiaChuyenSatPairs = [4][][2]int{
	{{1, 3}, {7, 3}},
	{{2, 6}, {4, 6}},
	{{7, 9}, {9, 9}},
	{{2, 0}, {8, 0}},
}

// tuPhePairs are the day pillars with Four Wasted, by season of the month branch.
var tuPhePairs = [4][][2]int{
	{{6, 8}, {7, 9}},
	{{8, 0}, {9, 11}},
	{{0, 2}, {1, 3}},
	{{2, 6}, {3, 5}},
}

// thienXichQuyPairs are the day pillars with Heavenly Pardon, by season of the month branch.
var thienXichQuyPairs = [4][][2]int{
	{{4, 2}},
	{{0, 6}},
	{{4, 8}},
	{{0, 0}},
}
