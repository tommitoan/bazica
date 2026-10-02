package bazica

import (
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/tommitoan/bazica/model"
)

func TestGetBaziChart(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	birth := time.Date(1995, 6, 8, 22, 5, 0, 0, loc)

	chart, err := GetBaziChart(birth, loc, 1)
	if err != nil {
		t.Fatalf("GetBaziChart() error = %v", err)
	}

	pillars := chart.FourPillar
	for _, tc := range []struct {
		name, stem, branch, lifeCycle string
		gotStem, gotBranch, gotCycle  string
	}{
		{"year", "yi", "hai", model.LC7, pillars.YearPillar.HeavenlyStem.Spelling, pillars.YearPillar.EarthlyBranch.Spelling, pillars.YearPillar.LifeCycle},
		{"month", "ren", "wu", model.LC2, pillars.MonthPillar.HeavenlyStem.Spelling, pillars.MonthPillar.EarthlyBranch.Spelling, pillars.MonthPillar.LifeCycle},
		{"day", "geng", "wu", model.LC2, pillars.DayPillar.HeavenlyStem.Spelling, pillars.DayPillar.EarthlyBranch.Spelling, pillars.DayPillar.LifeCycle},
		{"hour", "ding", "hai", model.LC7, pillars.HourPillar.HeavenlyStem.Spelling, pillars.HourPillar.EarthlyBranch.Spelling, pillars.HourPillar.LifeCycle},
	} {
		if tc.gotStem != tc.stem || tc.gotBranch != tc.branch {
			t.Errorf("%s pillar = %s/%s, want %s/%s", tc.name, tc.gotStem, tc.gotBranch, tc.stem, tc.branch)
		}
		if tc.gotCycle != tc.lifeCycle {
			t.Errorf("%s life cycle = %q, want %q", tc.name, tc.gotCycle, tc.lifeCycle)
		}
	}

	// Yin Wood year and male: pillars run backwards from Ren Wu, and the
	// start age is counted from the 3562 minutes since Grain in Ear.
	luck := chart.LuckPillars.LuckPillars
	if len(luck) != 12 {
		t.Fatalf("got %d luck pillars, want 12", len(luck))
	}
	for i, want := range []struct{ stem, branch string }{
		{"ren", "wu"}, {"xin", "si"}, {"geng", "chen"}, {"ji", "mao"},
	} {
		got := luck[i]
		if got.HeavenlyStem.Spelling != want.stem || got.EarthlyBranch.Spelling != want.branch {
			t.Errorf("luck pillar %d = %s/%s, want %s/%s", i, got.HeavenlyStem.Spelling, got.EarthlyBranch.Spelling, want.stem, want.branch)
		}
	}
	if wantStart := time.Date(1996, 4, 3, 22, 5, 0, 0, loc); !luck[1].Time.Equal(wantStart) {
		t.Errorf("first luck pillar starts %v, want %v", luck[1].Time, wantStart)
	}
	if luck[1].YearStart != 1996 || luck[2].YearStart != 2006 {
		t.Errorf("luck pillar years = %d, %d, want 1996, 2006", luck[1].YearStart, luck[2].YearStart)
	}
}

func TestGetBaziChartJSONFieldNames(t *testing.T) {
	chart, err := GetBaziChart(time.Date(1995, 6, 8, 22, 5, 0, 0, time.UTC), time.UTC, 1)
	if err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(chart.FourPillar.HourPillar.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"hour":22,"minute":5}`; string(out) != want {
		t.Errorf("hour JSON = %s, want %s", out, want)
	}
}

// dateTime is read as wall-clock time in loc, whichever zone it was built in.
func TestGetBaziChartReadsDateTimeInLoc(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	local := time.Date(2024, 3, 11, 3, 0, 0, 0, loc)
	instantInUTC := local.UTC()

	want, err := GetBaziChart(local, loc, model.GenderMale)
	if err != nil {
		t.Fatal(err)
	}
	got, err := GetBaziChart(instantInUTC, loc, model.GenderMale)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := got.FourPillar.HourPillar.EarthlyBranch.Spelling, want.FourPillar.HourPillar.EarthlyBranch.Spelling; g != w {
		t.Errorf("hour branch = %s, want %s", g, w)
	}
	if g, w := got.FourPillar.DayPillar.HeavenlyStem.Spelling, want.FourPillar.DayPillar.HeavenlyStem.Spelling; g != w {
		t.Errorf("day stem = %s, want %s", g, w)
	}
	if got.FourPillar.HourPillar.Hour.Hour != 3 {
		t.Errorf("hour = %d, want 3", got.FourPillar.HourPillar.Hour.Hour)
	}

	// A nil location keeps the zone the time already carries.
	nilLoc, err := GetBaziChart(local, nil, model.GenderMale)
	if err != nil {
		t.Fatal(err)
	}
	if nilLoc.FourPillar.HourPillar.EarthlyBranch.Spelling != want.FourPillar.HourPillar.EarthlyBranch.Spelling {
		t.Error("nil loc should behave like the time's own location")
	}
}

func TestGetBaziChartRejectsInvalidGender(t *testing.T) {
	for _, g := range []int{-1, 2, 7} {
		chart, err := GetBaziChart(time.Date(2000, 6, 1, 12, 0, 0, 0, time.UTC), time.UTC, g)
		if !errors.Is(err, model.ErrInvalidGender) || chart != nil {
			t.Errorf("gender %d: got chart=%v err=%v, want ErrInvalidGender", g, chart, err)
		}
	}
}

// Every pillar of every chart must be a valid Sexagenary pair (stem and branch
// share polarity), and every field must be populated.
func TestGetBaziChartInvariants(t *testing.T) {
	step := 37*time.Hour + 11*time.Minute
	count := 0
	for d := time.Date(1900, 1, 2, 0, 0, 0, 0, time.UTC); d.Year() < 2099; d = d.Add(step) {
		chart, err := GetBaziChart(d, time.UTC, int(d.Unix()%2))
		if err != nil {
			t.Fatalf("%s: %v", d.Format(time.RFC3339), err)
		}
		fp := chart.FourPillar
		for name, p := range map[string]struct {
			stem   model.HeavenlyStem
			branch model.EarthlyBranch
			gz     model.GanZhi
			lc     string
		}{
			"year":  {fp.YearPillar.HeavenlyStem, fp.YearPillar.EarthlyBranch, fp.YearPillar.GanZhi, fp.YearPillar.LifeCycle},
			"month": {fp.MonthPillar.HeavenlyStem, fp.MonthPillar.EarthlyBranch, fp.MonthPillar.GanZhi, fp.MonthPillar.LifeCycle},
			"day":   {fp.DayPillar.HeavenlyStem, fp.DayPillar.EarthlyBranch, fp.DayPillar.GanZhi, fp.DayPillar.LifeCycle},
			"hour":  {fp.HourPillar.HeavenlyStem, fp.HourPillar.EarthlyBranch, fp.HourPillar.GanZhi, fp.HourPillar.LifeCycle},
		} {
			if p.stem.Value == 0 || p.branch.Value == 0 || p.stem.Value%2 != p.branch.Value%2 {
				t.Fatalf("%s %s pillar is not a valid pair: stem %d branch %d", d.Format(time.RFC3339), name, p.stem.Value, p.branch.Value)
			}
			if p.gz.Name == "" || p.lc == "" {
				t.Fatalf("%s %s pillar is missing GanZhi or life cycle", d.Format(time.RFC3339), name)
			}
		}
		if n := len(chart.LuckPillars.LuckPillars); n != 12 {
			t.Fatalf("%s: %d luck pillars", d.Format(time.RFC3339), n)
		}
		count++
	}
	if count < 15000 {
		t.Fatalf("only %d dates checked", count)
	}
}

func TestGetBaziChartErrors(t *testing.T) {
	tests := []struct {
		name    string
		date    time.Time
		wantErr error
	}{
		{"before supported range", time.Date(1899, 12, 31, 12, 0, 0, 0, time.UTC), model.ErrDateOutOfRange},
		{"after supported range", time.Date(2100, 1, 1, 12, 0, 0, 0, time.UTC), model.ErrDateOutOfRange},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart, err := GetBaziChart(tt.date, time.UTC, 1)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetBaziChart() error = %v, want %v", err, tt.wantErr)
			}
			if chart != nil {
				t.Errorf("GetBaziChart() returned a chart alongside an error")
			}
		})
	}
}

// Calendar tables are loaded once and shared, so concurrent callers (for
// example HTTP handlers) must not race. Run with -race to exercise this.
func TestGetBaziChartConcurrent(t *testing.T) {
	const workers = 8
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				date := time.Date(1980+w, time.Month(i%12+1), i%28+1, i%24, 0, 0, 0, time.UTC)
				if _, err := GetBaziChart(date, time.UTC, w%2); err != nil {
					errs <- err
					return
				}
			}
		}(w)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}
