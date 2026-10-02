package bazica_test

import (
	"fmt"
	"time"

	"github.com/tommitoan/bazica"
	"github.com/tommitoan/bazica/model"
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
