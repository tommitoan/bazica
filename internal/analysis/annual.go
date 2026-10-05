package analysis

import (
	"fmt"

	"github.com/tommitoan/bazica/v2/internal/fourpillars"
	"github.com/tommitoan/bazica/v2/internal/utils"
	"github.com/tommitoan/bazica/v2/model"
)

// maxAnnualYear is the last Gregorian year of the yearly table. The first is the
// first year a birth date is supported (fourpillars.SupportedYears). The last is
// not tied to the calendar data: a yearly pillar is plain sixty-year arithmetic
// and needs no table, so the table may run past the last birth date. The bound
// keeps the answer finite (at most about eight thousand rows).
const maxAnnualYear = 9999

// AnnualPillars returns the pillars of count consecutive years from fromYear.
// The first year may not precede the year of the chart's year pillar, because
// nominal age would not be positive, and the last year may not pass 9999.
func AnnualPillars(chart *model.BaziChart, fromYear, count int) (*model.AnnualPillars, error) {
	if chart == nil {
		return nil, fmt.Errorf("%w: chart is required", model.ErrInvalidYearRange)
	}
	pillars, err := natalPillarsOf(chart)
	if err != nil {
		return nil, err
	}
	minAnnualYear, _, err := fourpillars.SupportedYears()
	if err != nil {
		return nil, err
	}
	baseYear := yearPillarYear(chart.FourPillar.YearPillar)
	last := fromYear + count - 1
	if count < 1 || fromYear < minAnnualYear || fromYear < baseYear || last > maxAnnualYear {
		return nil, model.ErrInvalidYearRange
	}

	dayStem := pillars[2].stem.Value
	natal := natalCells(pillars)
	rows := make([]model.AnnualPillar, 0, count)
	for year := fromYear; year <= last; year++ {
		// 1984 was a Jia-Zi year, the start of a sixty-year cycle.
		stem, branch := cycleStemBranch(year)
		s, b := utils.CalculateHeavenlyStem(stem), utils.CalculateEarthlyBranch(branch)
		ganZhi := utils.GanZhiOf(s, b)
		nayin, ok := nayinTerms[ganZhi.Name]
		if !ok {
			return nil, fmt.Errorf("%w: no Nayin for %s %s", errInvalidPillar, s.Name, b.Name)
		}
		rows = append(rows, model.AnnualPillar{
			Year:                 year,
			NominalAge:           year - baseYear + 1,
			HeavenlyStem:         s,
			EarthlyBranch:        b,
			GanZhi:               ganZhi,
			Nayin:                nayin,
			TenGod:               tenGod(dayStem, stem),
			StemStageAtOwnBranch: stage(stem, branch),
			LuckPillarNumber:     luckPillarOf(chart, year),
			HeavenEarthClash:     heavenEarthClash(cell{stem, branch}, natal),
		})
	}
	return &model.AnnualPillars{FromYear: fromYear, Count: count, AnnualPillars: rows}, nil
}

// cycleStemBranch returns the stem and branch values of a Gregorian year. 1984
// was a Jia-Zi year, the start of a sixty-year cycle.
func cycleStemBranch(year int) (stem, branch int) {
	position := ((year-1984)%60 + 60) % 60
	return stemValue(position), branchValue(position)
}

// yearPillarYear returns the Gregorian year the year pillar stands for. The
// model's YearPillar.Year is the civil year of the birth, which is one too high
// for a birth before Lichun, so the pillar's own stem and branch
// decide between that year and the one before.
func yearPillarYear(yp *model.YearPillar) int {
	for _, year := range []int{yp.Year, yp.Year - 1} {
		if stem, branch := cycleStemBranch(year); stem == yp.HeavenlyStem.Value && branch == yp.EarthlyBranch.Value {
			return year
		}
	}
	return yp.Year
}

// luckPillarOf returns the number of the luck pillar whose years contain year, or nil.
func luckPillarOf(chart *model.BaziChart, year int) *int {
	if chart.LuckPillars == nil {
		return nil
	}
	for _, lp := range chart.LuckPillars.LuckPillars {
		if lp != nil && lp.YearStart <= year && year <= lp.YearEnd {
			number := lp.Number
			return &number
		}
	}
	return nil
}
