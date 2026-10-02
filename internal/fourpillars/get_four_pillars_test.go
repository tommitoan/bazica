package fourpillars

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/tommitoan/bazica/model"
)

const dataPath = "../../"

func summarize(p *model.FourPillars) string {
	pillar := func(stem model.HeavenlyStem, branch model.EarthlyBranch, g model.GanZhi) string {
		return fmt.Sprintf("%s/%s %s", stem.Spelling, branch.Spelling, g.Name)
	}
	return fmt.Sprintf("Y[%s] M[%s] D[%s] H[%s]",
		pillar(p.YearPillar.HeavenlyStem, p.YearPillar.EarthlyBranch, p.YearPillar.GanZhi),
		pillar(p.MonthPillar.HeavenlyStem, p.MonthPillar.EarthlyBranch, p.MonthPillar.GanZhi),
		pillar(p.DayPillar.HeavenlyStem, p.DayPillar.EarthlyBranch, p.DayPillar.GanZhi),
		pillar(p.HourPillar.HeavenlyStem, p.HourPillar.EarthlyBranch, p.HourPillar.GanZhi))
}

func TestGetFourPillars(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			"late Rat hour rolls to next day",
			time.Date(1977, 7, 12, 23, 30, 0, 0, loc),
			"Y[ding/si Desert earth] M[ding/wei Sky water] D[xin/wei Road earth] H[wu/zi Lightning fire]",
		},
		{
			"mid-year birth",
			time.Date(1995, 6, 8, 22, 5, 0, 0, loc),
			"Y[yi/hai Volcanic fire] M[ren/wu Willow wood] D[geng/wu Road earth] H[ding/hai Roof tiles earth]",
		},
		{
			// Winter Solstice 2025-12-21 15:03 UTC; the last ten days of
			// December sit after it and previously caused an index panic.
			"after winter solstice",
			time.Date(2025, 12, 25, 12, 0, 0, 0, loc),
			"Y[yi/si Lamp fire] M[wu/zi Lightning fire] D[wu/chen Forest wood] H[wu/wu Sun fire]",
		},
		{
			"last day of the supported range",
			time.Date(2099, 12, 31, 12, 0, 0, 0, loc),
			"",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _, _, err := GetFourPillars(tc.date, loc, dataPath)
			if err != nil {
				t.Fatalf("GetFourPillars() error = %v", err)
			}
			if tc.want == "" {
				return
			}
			if s := summarize(got); s != tc.want {
				t.Errorf("GetFourPillars()\n got: %s\nwant: %s", s, tc.want)
			}
		})
	}
}

func TestGetFourPillarsErrors(t *testing.T) {
	loc := time.UTC
	tests := []struct {
		name    string
		date    time.Time
		path    string
		wantErr error
	}{
		{"before supported range", time.Date(1899, 6, 1, 12, 0, 0, 0, loc), dataPath, model.ErrDateOutOfRange},
		{"after supported range", time.Date(2100, 6, 1, 12, 0, 0, 0, loc), dataPath, model.ErrDateOutOfRange},
		{"first month of 2100", time.Date(2100, 1, 15, 12, 0, 0, 0, loc), dataPath, model.ErrDateOutOfRange},
		{"missing data directory", time.Date(2001, 1, 25, 12, 0, 0, 0, loc), "/nonexistent/", model.ErrDataUnavailable},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, err := GetFourPillars(tc.date, loc, tc.path)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("GetFourPillars() error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}
