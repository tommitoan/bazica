package fourpillars

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/tommitoan/bazica/model"
)

func spellings(p *model.FourPillars) string {
	return fmt.Sprintf("%s/%s %s/%s %s/%s %s/%s",
		p.YearPillar.HeavenlyStem.Spelling, p.YearPillar.EarthlyBranch.Spelling,
		p.MonthPillar.HeavenlyStem.Spelling, p.MonthPillar.EarthlyBranch.Spelling,
		p.DayPillar.HeavenlyStem.Spelling, p.DayPillar.EarthlyBranch.Spelling,
		p.HourPillar.HeavenlyStem.Spelling, p.HourPillar.EarthlyBranch.Spelling)
}

func ganZhiNames(p *model.FourPillars) string {
	return fmt.Sprintf("%s | %s | %s | %s",
		p.YearPillar.GanZhi.Name, p.MonthPillar.GanZhi.Name, p.DayPillar.GanZhi.Name, p.HourPillar.GanZhi.Name)
}

// The expected values come from an independent implementation: Julian-day
// arithmetic for the day pillar, the Five Tigers/Five Rats rules for the month
// and hour stems, and the latest "initial" solar term for the month branch.
// Year pillars follow the Lunar New Year, as documented in the README.
func TestGetFourPillarsGolden(t *testing.T) {
	tests := []struct {
		name string
		date string // UTC, "2006-01-02 15:04"
		want string // year month day hour as stem/branch
	}{
		{"start of supported range", "1900-01-01 12:00", "ji/hai bing/zi jia/xu geng/wu"},
		{"after Lunar New Year 1900", "1900-02-19 12:00", "geng/zi wu/yin gui/hai wu/wu"},
		{"after Lunar New Year and Lichun 1995", "1995-02-10 12:00", "yi/hai wu/yin ren/shen bing/wu"},
		{"shortly after Lunar New Year 1958", "1958-02-20 12:00", "wu/xu jia/yin wu/chen wu/wu"},
		{"after Lunar New Year 1985", "1985-02-25 12:00", "yi/chou wu/yin yi/wei ren/wu"},
		{"before Lunar New Year 1990", "1990-01-25 12:00", "ji/si ding/chou geng/yin ren/wu"},
		{"before Lunar New Year 2001", "2001-01-23 12:00", "geng/chen ji/chou bing/xu jia/wu"},
		{"early January belongs to the Rat month", "2024-01-03 12:00", "gui/mao jia/zi bing/yin jia/wu"},
		{"Ox month after Minor Cold", "2024-01-10 12:00", "gui/mao yi/chou gui/you wu/wu"},
		{"Ox month in a repaired data year", "2057-01-10 12:00", "bing/zi xin/chou ding/mao bing/wu"},
		{"Rat month in a repaired data year", "2057-12-29 12:00", "ding/chou ren/zi geng/shen ren/wu"},
		{"before Minor Cold in a repaired data year", "2041-01-03 12:00", "geng/shen wu/zi bing/shen jia/wu"},
		{"after winter solstice", "2025-12-25 12:00", "yi/si wu/zi wu/chen wu/wu"},
		{"last day of the year, late Rat hour", "2024-12-31 23:30", "jia/chen bing/zi geng/wu bing/zi"},
		{"hour boundary 00:59", "2024-03-10 00:59", "jia/chen ding/mao gui/you ren/zi"},
		{"hour boundary 01:00", "2024-03-10 01:00", "jia/chen ding/mao gui/you gui/chou"},
		{"hour boundary 22:59", "2024-03-10 22:59", "jia/chen ding/mao gui/you gui/hai"},
		{"hour boundary 23:00", "2024-03-10 23:00", "jia/chen ding/mao jia/xu jia/zi"},
		{"hour boundary 23:59", "2024-03-10 23:59", "jia/chen ding/mao jia/xu jia/zi"},
		{"two minutes before Minor Heat 2024", "2024-07-06 14:18", "jia/chen geng/wu xin/wei yi/wei"},
		{"two minutes after Minor Heat 2024", "2024-07-06 14:22", "jia/chen xin/wei xin/wei yi/wei"},
		{"Yin day master gui", "2020-03-01 12:00", "geng/zi wu/yin gui/mao wu/wu"},
		{"Yin day master yi", "2020-03-03 12:00", "geng/zi wu/yin yi/si ren/wu"},
		{"Yin day master ding", "2020-03-05 12:00", "geng/zi ji/mao ding/wei bing/wu"},
		{"Yin day master ji", "2020-03-07 12:00", "geng/zi ji/mao ji/you geng/wu"},
		{"Yin day master xin", "2020-03-09 12:00", "geng/zi ji/mao xin/hai jia/wu"},
		{"end of supported range", "2099-12-31 23:30", "ji/wei bing/zi gui/mao ren/zi"},
		// The cases below lie outside 1900-2099. They were derived with the independent
		// low-precision calculator used for the reference-page samples (Julian-day day
		// pillar, Sun longitude series for year and month) on days far from any boundary.
		{"first supported day", "1700-01-01 12:00", "ji/mao bing/zi bing/wu jia/wu"},
		{"after Lunar New Year 1700", "1700-03-15 12:00", "geng/chen ji/mao ji/wei geng/wu"},
		{"autumn 1750", "1750-09-20 12:00", "geng/wu yi/you geng/yin ren/wu"},
		{"summer 1800", "1800-06-15 12:00", "geng/shen ren/wu yi/hai ren/wu"},
		{"autumn 1850", "1850-09-20 12:00", "geng/xu yi/you jia/xu geng/wu"},
		{"mid December 1899", "1899-12-15 12:00", "ji/hai bing/zi ding/si bing/wu"},
		{"before Lunar New Year 2100", "2100-01-10 12:00", "ji/wei ding/chou ren/zi bing/wu"},
		{"after Lunar New Year 2100", "2100-03-15 12:00", "geng/shen ji/mao bing/chen jia/wu"},
		{"summer 2250", "2250-08-20 12:00", "geng/yin jia/shen geng/zi ren/wu"},
		{"summer 2399", "2399-06-10 12:00", "ji/wei geng/wu geng/yin ren/wu"},
		{"last day of the supported range, late Rat hour", "2399-12-31 23:30", "ji/wei bing/zi yi/mao bing/zi"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			date, err := time.ParseInLocation("2006-01-02 15:04", tc.date, time.UTC)
			if err != nil {
				t.Fatal(err)
			}
			got, _, _, err := GetFourPillars(date, time.UTC)
			if err != nil {
				t.Fatalf("GetFourPillars() error = %v", err)
			}
			if s := spellings(got); s != tc.want {
				t.Errorf("GetFourPillars()\n got: %s\nwant: %s", s, tc.want)
			}
		})
	}
}

func TestGetFourPillarsLocalClock(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	tests := []struct {
		name       string
		date       time.Time
		wantStems  string
		wantGanZhi string
	}{
		{
			"late Rat hour rolls to the next day",
			time.Date(1977, 7, 12, 23, 30, 0, 0, loc),
			"ding/si ding/wei xin/wei wu/zi",
			"Desert earth | Sky water | Road earth | Lightning fire",
		},
		{
			"evening birth",
			time.Date(1995, 6, 8, 22, 5, 0, 0, loc),
			"yi/hai ren/wu geng/wu ding/hai",
			"Volcanic fire | Willow wood | Road earth | Roof tiles earth",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _, _, err := GetFourPillars(tc.date, loc)
			if err != nil {
				t.Fatalf("GetFourPillars() error = %v", err)
			}
			if s := spellings(got); s != tc.wantStems {
				t.Errorf("pillars\n got: %s\nwant: %s", s, tc.wantStems)
			}
			if s := ganZhiNames(got); s != tc.wantGanZhi {
				t.Errorf("gan zhi\n got: %s\nwant: %s", s, tc.wantGanZhi)
			}
		})
	}
}

func TestGetFourPillarsErrors(t *testing.T) {
	loc := time.UTC
	tests := []struct {
		name    string
		date    time.Time
		wantErr error
	}{
		{"before supported range", time.Date(1699, 6, 1, 12, 0, 0, 0, loc), model.ErrDateOutOfRange},
		{"after supported range", time.Date(2400, 6, 1, 12, 0, 0, 0, loc), model.ErrDateOutOfRange},
		{"first month of 2400", time.Date(2400, 1, 15, 12, 0, 0, 0, loc), model.ErrDateOutOfRange},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, err := GetFourPillars(tc.date, loc)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("GetFourPillars() error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}
