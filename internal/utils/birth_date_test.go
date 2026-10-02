package utils

import (
	"testing"
	"time"
)

func TestBirthCalendarDate(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
		want string
	}{
		{"morning", time.Date(2024, 3, 10, 9, 0, 0, 0, time.UTC), "2024-03-10"},
		{"22:59 stays", time.Date(2024, 3, 10, 22, 59, 0, 0, time.UTC), "2024-03-10"},
		{"23:00 rolls over", time.Date(2024, 3, 10, 23, 0, 0, 0, time.UTC), "2024-03-11"},
		{"end of month", time.Date(2024, 2, 29, 23, 30, 0, 0, time.UTC), "2024-03-01"},
		{"end of year", time.Date(2024, 12, 31, 23, 59, 0, 0, time.UTC), "2025-01-01"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			y, m, d := BirthCalendarDate(tc.in)
			if got := time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Format("2006-01-02"); got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}
