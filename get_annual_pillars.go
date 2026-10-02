package bazica

import (
	"github.com/tommitoan/bazica/internal/analysis"
	"github.com/tommitoan/bazica/model"
)

// GetAnnualPillars returns the pillars of count consecutive Gregorian years
// starting at fromYear, each with its Ten God, life stage, nominal age, owning
// luck pillar and heaven-clash/earth-clash flag against the natal pillars.
//
// The nominal age counts the year of the chart's year pillar as 1, so a chart
// whose birth falls before the Lunar New Year starts at the previous year.
// fromYear may not precede that year, count must be positive, and the years
// must lie within 1900-2099; otherwise the error is model.ErrInvalidYearRange.
func GetAnnualPillars(chart *model.BaziChart, fromYear, count int) (*model.AnnualPillars, error) {
	return analysis.AnnualPillars(chart, fromYear, count)
}
