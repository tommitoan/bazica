package model

import "errors"

var (
	// ErrDateOutOfRange is returned when the requested date falls outside the
	// range covered by the bundled calendar data; bazica.SupportedYears returns
	// that range.
	ErrDateOutOfRange = errors.New("bazica: date is outside the supported range")

	// ErrDataUnavailable is returned when the bundled calendar data cannot be
	// read or decoded.
	ErrDataUnavailable = errors.New("bazica: calendar data unavailable")

	// ErrInvalidYearRange is returned by GetAnnualPillars when the chart is
	// missing, the count is not positive, or the years fall before the year
	// pillar's year, before the first supported birth year or past 9999.
	ErrInvalidYearRange = errors.New("bazica: annual pillar years must be positive in count, not before the year pillar or the first supported year, and not past 9999")

	// ErrInvalidGender is returned when the gender is neither GenderFemale nor GenderMale.
	ErrInvalidGender = errors.New("bazica: gender must be 0 (female) or 1 (male)")
)

// Gender values accepted by GetBaziChart.
const (
	GenderFemale = 0
	GenderMale   = 1
)
