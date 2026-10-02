package model

import "errors"

var (
	// ErrDateOutOfRange is returned when the requested date falls outside the
	// range covered by the bundled calendar data (1900-01-01 to 2099-12-31).
	ErrDateOutOfRange = errors.New("bazica: date is outside the supported range 1900-2099")

	// ErrDataUnavailable is returned when the bundled calendar data cannot be
	// read or decoded.
	ErrDataUnavailable = errors.New("bazica: calendar data unavailable")

	// ErrInvalidYearRange is returned by GetAnnualPillars when the chart is
	// missing, the count is not positive, or the years fall before the year
	// pillar's year or outside 1900-2099.
	ErrInvalidYearRange = errors.New("bazica: annual pillar years must be positive in count, not before the year pillar and within 1900-2099")

	// ErrInvalidGender is returned when the gender is neither GenderFemale nor GenderMale.
	ErrInvalidGender = errors.New("bazica: gender must be 0 (female) or 1 (male)")
)

// Gender values accepted by GetBaziChart.
const (
	GenderFemale = 0
	GenderMale   = 1
)
