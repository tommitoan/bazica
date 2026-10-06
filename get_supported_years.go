package bazica

import "github.com/tommitoan/bazica/v2/internal/fourpillars"

// SupportedYears returns the first and last calendar year of a birth date that
// GetBaziChart accepts: from 1 January of the first year to 31 December of the
// last. The range follows the embedded calendar tables, so callers that need it
// (a date picker, input validation) should read it here instead of repeating
// the numbers. The error wraps model.ErrDataUnavailable if the tables cannot be
// read.
func SupportedYears() (first, last int, err error) {
	return fourpillars.SupportedYears()
}
