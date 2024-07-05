package bazica

import (
	"fmt"
	"testing"
	"time"
)

func TestDetectSolarTerm(t *testing.T) {
	birthYear := 1900
	birthMonth := time.January
	birthDay := 2
	birthHour := 11
	birthMinute := 05
	location := "Asia/Ho_Chi_Minh"

	loc, err := time.LoadLocation(location)
	dateTime := time.Date(birthYear, birthMonth, birthDay, birthHour, birthMinute, 0, 0, loc)
	str, err := GetSolarTerm("", dateTime)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Printf("Solar term for %v: %s\n", dateTime, str)
}
