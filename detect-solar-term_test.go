package bazica

import (
	"fmt"
	"testing"
)

func TestDetectSolarTerm(t *testing.T) {
	input := "1945-03-06 11:32:29.605+00:00"
	str, err := DetectSolarTerm("", input)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Printf("Solar term for %s: %s\n", input, str)

}
