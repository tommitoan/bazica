package model

import "errors"

var (
	// ErrDateOutOfRange is returned when the requested date falls outside the
	// range covered by the bundled calendar data (1900-01-01 to 2099-12-31).
	ErrDateOutOfRange = errors.New("bazica: date is outside the supported range 1900-2099")

	// ErrDataUnavailable is returned when the bundled calendar data cannot be
	// read or decoded.
	ErrDataUnavailable = errors.New("bazica: calendar data unavailable")
)
