package bazica

import (
	"errors"
	"strings"
	"time"
)

func GetDateFormat(input string) (time.Time, error) {
	dateParts := strings.Split(input, "T") // Split date and time (if present)
	if len(dateParts) == 0 {
		return time.Time{}, errors.New("date part is empty")
	}
	dateStr := dateParts[0]
	solarDate, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return solarDate, nil
}
