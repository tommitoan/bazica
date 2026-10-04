package bazica

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"
)

// Every star in a chart carries its nature, nothing else does, and the legacy
// fields stay as they were (TestLegacyFieldsAreUnchanged covers the rest).
func TestChartStarsCarryTheirNature(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	births := []time.Time{
		time.Date(1995, 6, 8, 11, 1, 0, 0, loc),
		time.Date(1964, 11, 30, 23, 30, 0, 0, loc),
		time.Date(2036, 6, 20, 3, 5, 0, 0, loc),
	}
	starCode := regexp.MustCompile(`"code":"star\.[a-z_]+"`)
	starWithNature := regexp.MustCompile(`"code":"star\.[a-z_]+","en":"[^"]*","vi":"[^"]*","nature":"(auspicious|inauspicious|mixed)"`)
	for _, birth := range births {
		for gender := 0; gender <= 1; gender++ {
			chart, err := GetBaziChart(birth, loc, gender)
			if err != nil {
				t.Fatalf("GetBaziChart(%v, %d) error = %v", birth, gender, err)
			}
			data, err := json.Marshal(chart)
			if err != nil {
				t.Fatal(err)
			}
			stars := len(starCode.FindAll(data, -1))
			if stars == 0 {
				t.Fatalf("%v: the chart has no stars to check", birth)
			}
			if with := len(starWithNature.FindAll(data, -1)); with != stars {
				t.Errorf("%v gender %d: %d stars, %d with a nature", birth, gender, stars, with)
			}
			if all := len(regexp.MustCompile(`"nature":`).FindAll(data, -1)); all != stars {
				t.Errorf("%v gender %d: %d nature keys for %d stars", birth, gender, all, stars)
			}
		}
	}
}
