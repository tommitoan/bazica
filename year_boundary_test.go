package bazica

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/model"
)

func pillarChars(s model.HeavenlyStem, b model.EarthlyBranch) string {
	return s.Character + b.Character
}

// luckForward reports whether the first luck pillar follows the month pillar in the sixty-cycle.
func luckForward(t *testing.T, chart *model.BaziChart) bool {
	t.Helper()
	const stems, branches = "甲乙丙丁戊己庚辛壬癸", "子丑寅卯辰巳午未申酉戌亥"
	index := func(gz string) int {
		runes := []rune(gz)
		var s, b int = -1, -1
		for i, r := range stems {
			if r == runes[0] {
				s = i / 3
			}
		}
		for i, r := range branches {
			if r == runes[1] {
				b = i / 3
			}
		}
		for n := 0; n < 60; n++ {
			if n%10 == s && n%12 == b {
				return n
			}
		}
		t.Fatalf("no position for %s", gz)
		return 0
	}
	month := pillarChars(chart.FourPillar.MonthPillar.HeavenlyStem, chart.FourPillar.MonthPillar.EarthlyBranch)
	first := chart.LuckPillars.LuckPillars[1]
	return (index(pillarChars(first.HeavenlyStem, first.EarthlyBranch))-index(month)+60)%60 == 1
}

// Days between the Lunar New Year and Lichun (in either order) of 1901-2099 give the
// pillars, the first luck pillar and the luck direction of a second implementation.
// The fixture was generated with lunar-javascript; see its note.
func TestWindowDaysMatchSecondImplementation(t *testing.T) {
	raw, err := os.ReadFile("testdata/calendar/lunarjs-window-1901-2099.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Cases []struct {
			Date    string `json:"date"`
			Gender  int    `json:"gender"`
			Year    string `json:"year"`
			Month   string `json:"month"`
			Day     string `json:"day"`
			Hour    string `json:"hour"`
			Luck1   string `json:"luck1"`
			Forward bool   `json:"forward"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(raw, &fixture); err != nil {
		t.Fatal(err)
	}
	if len(fixture.Cases) < 4000 {
		t.Fatalf("fixture holds %d cases, expected about 4000", len(fixture.Cases))
	}
	loc := time.FixedZone("UTC+8", 8*3600)
	for _, c := range fixture.Cases {
		day, err := time.ParseInLocation("2006-01-02", c.Date, loc)
		if err != nil {
			t.Fatal(err)
		}
		birth := time.Date(day.Year(), day.Month(), day.Day(), 12, 0, 0, 0, loc)
		chart, err := GetBaziChart(birth, loc, c.Gender)
		if err != nil {
			t.Fatalf("%s: %v", c.Date, err)
		}
		p := chart.FourPillar
		got := [...]string{
			pillarChars(p.YearPillar.HeavenlyStem, p.YearPillar.EarthlyBranch),
			pillarChars(p.MonthPillar.HeavenlyStem, p.MonthPillar.EarthlyBranch),
			pillarChars(p.DayPillar.HeavenlyStem, p.DayPillar.EarthlyBranch),
			pillarChars(p.HourPillar.HeavenlyStem, p.HourPillar.EarthlyBranch),
		}
		want := [...]string{c.Year, c.Month, c.Day, c.Hour}
		if got != want {
			t.Fatalf("%s gender %d: pillars %v, second implementation %v", c.Date, c.Gender, got, want)
		}
		first := chart.LuckPillars.LuckPillars[1]
		if luck := pillarChars(first.HeavenlyStem, first.EarthlyBranch); luck != c.Luck1 {
			t.Fatalf("%s gender %d: first luck pillar %s, second implementation %s", c.Date, c.Gender, luck, c.Luck1)
		}
		if forward := luckForward(t, chart); forward != c.Forward {
			t.Fatalf("%s gender %d: luck forward %v, second implementation %v", c.Date, c.Gender, forward, c.Forward)
		}
	}
}

// Ten days between the Lunar New Year and Lichun, named so a failure reads well. Years
// 1955-2018, noon at UTC+8; the values agree with the fixture above.
func TestWindowGoldens(t *testing.T) {
	tests := []struct {
		date               string
		gender             int
		year, month, luck1 string
		forward            bool
	}{
		{"1955-01-29", 0, "甲午", "丁丑", "丙子", false},
		{"1966-01-28", 1, "乙巳", "己丑", "戊子", false},
		{"1982-01-30", 0, "辛酉", "辛丑", "壬寅", true},
		{"1990-01-31", 1, "己巳", "丁丑", "丙子", false},
		{"2017-01-31", 0, "丙申", "辛丑", "庚子", false},
		{"1956-02-08", 1, "丙申", "庚寅", "辛卯", true},
		{"1958-02-11", 0, "戊戌", "甲寅", "癸丑", false},
		{"1980-02-10", 1, "庚申", "戊寅", "己卯", true},
		{"1994-02-07", 0, "甲戌", "丙寅", "乙丑", false},
		{"2018-02-10", 1, "戊戌", "甲寅", "乙卯", true},
	}
	loc := time.FixedZone("UTC+8", 8*3600)
	for _, tc := range tests {
		t.Run(tc.date, func(t *testing.T) {
			day, _ := time.ParseInLocation("2006-01-02", tc.date, loc)
			chart, err := GetBaziChart(time.Date(day.Year(), day.Month(), day.Day(), 12, 0, 0, 0, loc), loc, tc.gender)
			if err != nil {
				t.Fatal(err)
			}
			p := chart.FourPillar
			if got := pillarChars(p.YearPillar.HeavenlyStem, p.YearPillar.EarthlyBranch); got != tc.year {
				t.Errorf("year pillar %s, want %s", got, tc.year)
			}
			if got := pillarChars(p.MonthPillar.HeavenlyStem, p.MonthPillar.EarthlyBranch); got != tc.month {
				t.Errorf("month pillar %s, want %s", got, tc.month)
			}
			first := chart.LuckPillars.LuckPillars[1]
			if got := pillarChars(first.HeavenlyStem, first.EarthlyBranch); got != tc.luck1 {
				t.Errorf("first luck pillar %s, want %s", got, tc.luck1)
			}
			if got := luckForward(t, chart); got != tc.forward {
				t.Errorf("luck forward %v, want %v", got, tc.forward)
			}
		})
	}
}

// Charts that fall outside the two windows are unchanged from v1.5.0. The fixture holds the
// SHA-256 of each chart's canonical JSON as v1.5.0 produced it; fields added since then are
// not part of the comparison (see canonicalChartJSON).
func TestChartsOutsideTheWindowsAreUnchanged(t *testing.T) {
	raw, err := os.ReadFile("testdata/calendar/v150-regression-outside-windows.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		DateTime string `json:"dateTime"`
		Offset   int    `json:"offset"`
		Gender   int    `json:"gender"`
		Pillars  string `json:"pillars"`
		SHA256   string `json:"sha256"`
	}
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) < 600 {
		t.Fatalf("fixture holds %d cases, expected about 650", len(cases))
	}
	for _, c := range cases {
		loc := time.FixedZone("z", c.Offset*3600)
		birth, err := time.ParseInLocation("2006-01-02T15:04", c.DateTime, loc)
		if err != nil {
			t.Fatal(err)
		}
		chart, err := GetBaziChart(birth, loc, c.Gender)
		if err != nil {
			t.Fatalf("%s: %v", c.DateTime, err)
		}
		if got := canonicalChartJSON(t, chart); got != c.SHA256 {
			t.Fatalf("%s (offset %d, gender %d, pillars %s): chart differs from v1.5.0", c.DateTime, c.Offset, c.Gender, c.Pillars)
		}
	}
}

// canonicalChartJSON hashes the chart as generic JSON with sorted keys, leaving out blocks
// that were added after the fixture was made.
func canonicalChartJSON(t *testing.T, chart *model.BaziChart) string {
	t.Helper()
	b, err := json.Marshal(chart)
	if err != nil {
		t.Fatal(err)
	}
	var generic map[string]any
	if err := json.Unmarshal(b, &generic); err != nil {
		t.Fatal(err)
	}
	if analysis, ok := generic["analysis"].(map[string]any); ok {
		delete(analysis, "year_reference")
	}
	canon, err := json.Marshal(generic)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(canon)
	return hex.EncodeToString(sum[:])
}

// The reference block is display-only. Only the code that builds it may name it as an
// identifier (comments are ignored); any other non-test file that does would be
// calculating from it.
func TestNoCalculationReadsTheYearReference(t *testing.T) {
	allowed := map[string]bool{
		"model/model_year_reference.go":          true,
		"model/model_analysis.go":                true,
		"internal/fourpillars/year_reference.go": true,
		"internal/analysis/analyze.go":           true,
	}
	names := map[string]bool{"YearReference": true, "LunarYearReference": true, "LichunReference": true}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "testdata" || d.Name() == "tools") {
			return fs.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || allowed[filepath.ToSlash(path)] {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if ident, ok := n.(*ast.Ident); ok && names[ident.Name] {
				t.Errorf("%s uses %s; the year reference is display-only and may not feed a calculation", path, ident.Name)
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
