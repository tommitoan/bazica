package bazica_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/tommitoan/bazica/v2"
	"github.com/tommitoan/bazica/v2/model"
)

func ExampleGetBaziChart() {
	loc := time.FixedZone("UTC+7", 7*60*60)
	birth := time.Date(1995, time.June, 8, 22, 5, 0, 0, loc)

	chart, err := bazica.GetBaziChart(birth, loc, model.GenderMale)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := chart.FourPillar
	fmt.Println("year: ", p.YearPillar.HeavenlyStem.Name, p.YearPillar.EarthlyBranch.Name)
	fmt.Println("month:", p.MonthPillar.HeavenlyStem.Name, p.MonthPillar.EarthlyBranch.Name)
	fmt.Println("day:  ", p.DayPillar.HeavenlyStem.Name, p.DayPillar.EarthlyBranch.Name)
	fmt.Println("hour: ", p.HourPillar.HeavenlyStem.Name, p.HourPillar.EarthlyBranch.Name)
	// Output:
	// year:  Yin Wood Pig
	// month: Yang Water Horse
	// day:   Yang Metal Horse
	// hour:  Yin Fire Pig
}

func ExampleGetBaziChart_analysis() {
	loc := time.FixedZone("UTC+7", 7*60*60)
	birth := time.Date(1995, time.June, 8, 22, 5, 0, 0, loc)

	chart, err := bazica.GetBaziChart(birth, loc, model.GenderMale)
	if err != nil {
		fmt.Println(err)
		return
	}

	a := chart.Analysis
	fmt.Println("day master:  ", a.DayMaster.HeavenlyStem.Name, "-", a.DayMaster.Element.EN, "/", a.DayMaster.Element.VI)
	fmt.Println("void:        ", a.VoidBranches[0].Name, a.VoidBranches[1].Name)
	fmt.Println("life palace: ", a.LifePalace.HeavenlyStem.Name, a.LifePalace.EarthlyBranch.Name, "-", a.LifePalace.Nayin.VI)

	year := chart.FourPillar.YearPillar.Analysis
	fmt.Println("year ten god:", year.TenGod.EN, "/", year.TenGod.VI)
	var stars []string
	for _, star := range year.Stars {
		stars = append(stars, star.EN)
	}
	fmt.Println("year stars:  ", strings.Join(stars, ", "))
	// Output:
	// day master:   Yang Metal - Metal / Kim
	// void:         Dog Pig
	// life palace:  Yang Fire Rat - Giản Hạ Thủy
	// year ten god: Direct Wealth / Chính Tài
	// year stars:   National Seal, Literary Star, Robbery Sha, Heavenly Virtue Noble, Void
}

func ExampleGetAnnualPillars() {
	loc := time.FixedZone("UTC+7", 7*60*60)
	birth := time.Date(1995, time.June, 8, 22, 5, 0, 0, loc)

	chart, err := bazica.GetBaziChart(birth, loc, model.GenderMale)
	if err != nil {
		fmt.Println(err)
		return
	}
	annual, err := bazica.GetAnnualPillars(chart, 2026, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, y := range annual.AnnualPillars {
		fmt.Printf("%d age %d: %s %s, %s, luck pillar %d\n",
			y.Year, y.NominalAge, y.HeavenlyStem.Name, y.EarthlyBranch.Name, y.TenGod.EN, *y.LuckPillarNumber)
	}
	// Output:
	// 2026 age 32: Yang Fire Horse, Seven Killings, luck pillar 4
	// 2027 age 33: Yin Fire Goat, Direct Officer, luck pillar 4
}
