package bazica

import (
	"fmt"
	"testing"
)

func TestGetFourPillarChart(t *testing.T) {
	input := "1995-06-08  01:32:29.605+00:00"
	anything, err := GetFourPillarChart("", input)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(anything)
}
