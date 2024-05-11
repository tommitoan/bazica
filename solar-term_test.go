package bazica

import (
	"fmt"
	"testing"
)

func TestGetSolarTermData(t *testing.T) {
	result, _ := GetSolarTermsByYear("2000", "")

	fmt.Println(result)
	fmt.Println(result.ColdDew)
	fmt.Println(result.AwakeningOfInsects)
	fmt.Println(result.EndOfHeat)
	fmt.Println(result.GrainBuds)

}
