package fourpillars

import (
	"testing"
	"time"
	_ "time/tzdata" // keep the zone tests independent of the host's zoneinfo
)

// julianDay is an independent reference for the day count (Fliegel-Van Flandern).
func julianDay(y, m, d int) int {
	a := (14 - m) / 12
	yy := y + 4800 - a
	mm := m + 12*a - 3
	return d + (153*mm+2)/5 + 365*yy + yy/4 - yy/100 + yy/400 - 32045
}

// refDayPillar returns the stem value (Jia = 1) and the branch value (Tiger = 1)
// of the day a birth belongs to; from 23:00 that is the next calendar day.
func refDayPillar(t time.Time) (stem, branch int) {
	y, m, d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Date()
	if t.Hour() >= 23 {
		y, m, d = time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC).Date()
	}
	i := (julianDay(y, int(m), d) + 49) % 60 // 0 is Jia Zi
	return i%10 + 1, (i%12+10)%12 + 1
}

// Zones whose offset in 1900 differs from today's, or that observe daylight
// saving time, used to shift the day pillar for births around midnight.
func TestGetDayPillarAroundMidnightInOtherZones(t *testing.T) {
	zones := []string{"Asia/Ho_Chi_Minh", "America/New_York", "Europe/Paris", "Australia/Sydney"}
	hours := []int{0, 1, 22, 23}
	for _, name := range zones {
		loc, err := time.LoadLocation(name)
		if err != nil {
			t.Fatal(err)
		}
		t.Run(name, func(t *testing.T) {
			checked := 0
			for d := time.Date(1950, 1, 1, 0, 0, 0, 0, loc); d.Year() < 2030; d = d.AddDate(0, 0, 3) {
				for _, h := range hours {
					date := time.Date(d.Year(), d.Month(), d.Day(), h, 30, 0, 0, loc)
					if date.Hour() != h { // the hour does not exist on a DST change day
						continue
					}
					got, err := GetDayPillar(date)
					if err != nil {
						t.Fatal(err)
					}
					wantStem, wantBranch := refDayPillar(date)
					if got.HeavenlyStem.Value != wantStem || got.EarthlyBranch.Value != wantBranch {
						t.Fatalf("%s: got stem %d branch %d, want stem %d branch %d",
							date.Format("2006-01-02 15:04 MST"), got.HeavenlyStem.Value, got.EarthlyBranch.Value, wantStem, wantBranch)
					}
					checked++
				}
			}
			if checked == 0 {
				t.Fatal("no dates checked")
			}
		})
	}
}
