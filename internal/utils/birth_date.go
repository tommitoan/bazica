package utils

import "time"

// RatHourStart is the clock hour from which a birth counts toward the next day.
const RatHourStart = 23

// BirthCalendarDate returns the calendar date a birth belongs to. The Rat hour
// begins at 23:00, so from that hour on the date is the following day. The
// result uses the wall-clock date of t's own location, which keeps it correct
// across daylight saving and historical offset changes.
func BirthCalendarDate(t time.Time) (year int, month time.Month, day int) {
	year, month, day = t.Date()
	if t.Hour() >= RatHourStart {
		return time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC).Date()
	}
	return year, month, day
}
