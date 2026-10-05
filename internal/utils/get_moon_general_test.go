package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/tommitoan/bazica/model"
)

func TestGetMoonGeneralSwitchesOnThePrincipalTermDay(t *testing.T) {
	ict := time.FixedZone("ICT", 7*3600)
	tests := []struct {
		name string
		date time.Time
		want int
	}{
		// The Summer Solstice of 2022 fell on 21 June at 16:13 ICT: that whole day already has the new general.
		{"before the solstice day", time.Date(2022, 6, 20, 23, 59, 0, 0, ict), 8},
		{"morning of the solstice day", time.Date(2022, 6, 21, 3, 0, 0, 0, ict), 7},
		// The Winter Solstice of 2023 fell on 22 December at 10:27 ICT.
		{"before the winter solstice day", time.Date(2023, 12, 21, 12, 0, 0, 0, ict), 2},
		{"early on the winter solstice day", time.Date(2023, 12, 22, 3, 30, 0, 0, ict), 1},
		// Major Cold (Rat) follows the Winter Solstice (Ox).
		{"after major cold 2006", time.Date(2006, 1, 20, 20, 0, 0, 0, ict), 0},
		{"just before major cold 2006", time.Date(2006, 1, 19, 12, 0, 0, 0, ict), 1},
		// The first and last supported days still find a preceding term in the previous year's table.
		{"first supported day", time.Date(1700, 1, 1, 12, 0, 0, 0, ict), 1},
		{"last supported day", time.Date(2399, 12, 31, 12, 0, 0, 0, ict), 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetMoonGeneral(tc.date)
			if err != nil {
				t.Fatalf("GetMoonGeneral() error = %v", err)
			}
			if got != tc.want {
				t.Errorf("GetMoonGeneral(%v) = %d, want %d", tc.date, got, tc.want)
			}
		})
	}
}

func TestGetMoonGeneralOutsideTheTables(t *testing.T) {
	for _, date := range []time.Time{time.Date(1699, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2401, 6, 1, 0, 0, 0, 0, time.UTC)} {
		if _, err := GetMoonGeneral(date); !errors.Is(err, model.ErrDateOutOfRange) {
			t.Errorf("GetMoonGeneral(%v) error = %v, want ErrDateOutOfRange", date, err)
		}
	}
}
