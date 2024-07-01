package bazica

import (
	"fmt"
	"testing"
)

func TestGetYearPillar(t *testing.T) {
	input := "1945-03-06  00:32:29.605+00:00"
	str, err := GetYearPillar("", input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)

}
