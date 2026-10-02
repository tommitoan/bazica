// Package analysis derives Ba-zi facts that do not depend on the calendar:
// Ten Gods, hidden stems, life stages, void branches, auxiliary pillars,
// element counts and the heaven-clash/earth-clash flag. Every function is pure
// and works on stem and branch values of the model.
package analysis

import (
	"errors"
	"fmt"

	"github.com/tommitoan/bazica/internal/utils"
	"github.com/tommitoan/bazica/model"
)

// errInvalidPillar is returned for a pillar whose stem, branch or Nayin is not valid.
var errInvalidPillar = errors.New("bazica: invalid pillar for analysis")

// natalPillar is the part of a pillar the analysis reads.
type natalPillar struct {
	stem   model.HeavenlyStem
	branch model.EarthlyBranch
	ganZhi model.GanZhi
}

// Attach fills the Analysis fields of the four natal pillars, of the luck
// pillars (when present) and of the chart. The legacy fields of chart are left untouched.
func Attach(chart *model.BaziChart) error {
	pillars, err := natalPillarsOf(chart)
	if err != nil {
		return err
	}

	analyses, err := analyzePillars(pillars)
	if err != nil {
		return err
	}
	fp := chart.FourPillar
	fp.YearPillar.Analysis = analyses[0]
	fp.MonthPillar.Analysis = analyses[1]
	fp.DayPillar.Analysis = analyses[2]
	fp.HourPillar.Analysis = analyses[3]

	chartAnalysis, err := analyzeChart(pillars)
	if err != nil {
		return err
	}
	chart.Analysis = chartAnalysis
	return attachLuck(chart, pillars)
}

// natalPillarsOf reads and validates the four natal pillars of a chart.
func natalPillarsOf(chart *model.BaziChart) ([4]natalPillar, error) {
	var pillars [4]natalPillar
	fp := chart.FourPillar
	if fp == nil || fp.YearPillar == nil || fp.MonthPillar == nil || fp.DayPillar == nil || fp.HourPillar == nil {
		return pillars, fmt.Errorf("%w: four pillars are required", errInvalidPillar)
	}
	pillars = [4]natalPillar{
		{fp.YearPillar.HeavenlyStem, fp.YearPillar.EarthlyBranch, fp.YearPillar.GanZhi},
		{fp.MonthPillar.HeavenlyStem, fp.MonthPillar.EarthlyBranch, fp.MonthPillar.GanZhi},
		{fp.DayPillar.HeavenlyStem, fp.DayPillar.EarthlyBranch, fp.DayPillar.GanZhi},
		{fp.HourPillar.HeavenlyStem, fp.HourPillar.EarthlyBranch, fp.HourPillar.GanZhi},
	}
	for _, p := range pillars {
		if !validStem(p.stem.Value) || !validBranch(p.branch.Value) {
			return pillars, fmt.Errorf("%w: stem %d, branch %d", errInvalidPillar, p.stem.Value, p.branch.Value)
		}
	}
	return pillars, nil
}

// natalCells returns the four natal pillars as clash-rule cells.
func natalCells(pillars [4]natalPillar) []cell {
	cells := make([]cell, len(pillars))
	for i, p := range pillars {
		cells[i] = cell{p.stem.Value, p.branch.Value}
	}
	return cells
}

func analyzePillars(pillars [4]natalPillar) ([4]*model.PillarAnalysis, error) {
	var out [4]*model.PillarAnalysis
	dayStem, dayBranch := pillars[2].stem.Value, pillars[2].branch.Value
	monthBranch := pillars[1].branch.Value
	void := voidBranches(dayStem, dayBranch)

	cells := natalCells(pillars)

	for i, p := range pillars {
		nayin, ok := nayinTerms[p.ganZhi.Name]
		if !ok {
			return out, fmt.Errorf("%w: no Nayin for %q", errInvalidPillar, p.ganZhi.Name)
		}
		a := &model.PillarAnalysis{
			IsDayMaster:            i == 2,
			StemElement:            elementTerms[stemElement(p.stem.Value)],
			BranchElement:          elementTerms[branchElement(p.branch.Value)],
			Nayin:                  nayin,
			StemStageAtOwnBranch:   stage(p.stem.Value, p.branch.Value),
			StemStageAtMonthBranch: stage(p.stem.Value, monthBranch),
			DayMasterStage:         stage(dayStem, p.branch.Value),
			HiddenStems:            hiddenStemList(dayStem, p.branch.Value, monthBranch),
			IsVoid:                 i != 2 && (p.branch.Value == void[0] || p.branch.Value == void[1]),
			HeavenEarthClash:       heavenEarthClash(cells[i], othersOf(cells, i)),
			Stars:                  []model.LocalizedTerm{},
		}
		if i != 2 {
			god := tenGod(dayStem, p.stem.Value)
			a.TenGod = &god
		}
		out[i] = a
	}
	return out, nil
}

// othersOf returns every cell except the one at index skip.
func othersOf(cells []cell, skip int) []cell {
	others := make([]cell, 0, len(cells)-1)
	for i, c := range cells {
		if i != skip {
			others = append(others, c)
		}
	}
	return others
}

func hiddenStemList(dayStem, branch, monthBranch int) []model.HiddenStem {
	values := hiddenStemsOf(branch)
	list := make([]model.HiddenStem, 0, len(values))
	for _, h := range values {
		list = append(list, model.HiddenStem{
			HeavenlyStem:       utils.CalculateHeavenlyStem(h),
			TenGod:             tenGod(dayStem, h),
			StageAtMonthBranch: stage(h, monthBranch),
		})
	}
	return list
}

func analyzeChart(pillars [4]natalPillar) (*model.ChartAnalysis, error) {
	year, month, day, hour := pillars[0], pillars[1], pillars[2], pillars[3]
	void := voidBranches(day.stem.Value, day.branch.Value)

	tnStem, tnBranch := thaiNguyen(month.stem.Value, month.branch.Value)
	ttStem, ttBranch := thaiTuc(day.stem.Value, day.branch.Value)
	lpStem, lpBranch := lifePalace(year.stem.Value, month.branch.Value, hour.branch.Value)
	thaiNguyenPillar, err := auxiliaryPillar(tnStem, tnBranch)
	if err != nil {
		return nil, err
	}
	thaiTucPillar, err := auxiliaryPillar(ttStem, ttBranch)
	if err != nil {
		return nil, err
	}
	lifePalacePillar, err := auxiliaryPillar(lpStem, lpBranch)
	if err != nil {
		return nil, err
	}

	var counts model.ElementCounts
	for _, p := range pillars {
		addElement(&counts.Stems, stemElement(p.stem.Value))
		addElement(&counts.Branches, branchElement(p.branch.Value))
		for _, h := range hiddenStemsOf(p.branch.Value) {
			addElement(&counts.HiddenStems, stemElement(h))
		}
	}

	return &model.ChartAnalysis{
		DayMaster: model.DayMaster{
			HeavenlyStem: day.stem,
			Element:      elementTerms[stemElement(day.stem.Value)],
			Polarity:     polarityTerm(day.stem.Value),
		},
		VoidBranches:  []model.EarthlyBranch{utils.CalculateEarthlyBranch(void[0]), utils.CalculateEarthlyBranch(void[1])},
		ThaiNguyen:    thaiNguyenPillar,
		ThaiTuc:       thaiTucPillar,
		LifePalace:    lifePalacePillar,
		ElementCounts: counts,
	}, nil
}

func auxiliaryPillar(stem, branch int) (model.AuxiliaryPillar, error) {
	s, b := utils.CalculateHeavenlyStem(stem), utils.CalculateEarthlyBranch(branch)
	ganZhi := utils.GanZhiOf(s, b)
	nayin, ok := nayinTerms[ganZhi.Name]
	if !ok {
		return model.AuxiliaryPillar{}, fmt.Errorf("%w: no Nayin for %s %s", errInvalidPillar, s.Name, b.Name)
	}
	return model.AuxiliaryPillar{HeavenlyStem: s, EarthlyBranch: b, GanZhi: ganZhi, Nayin: nayin}, nil
}
