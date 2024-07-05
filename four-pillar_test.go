package bazica

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestGetFourPillarChart(t *testing.T) {
	// Test Cases: Each case defines a birthdate, expected chart, and (optionally) a description.
	testCases := []struct {
		loc                            string
		year, month, day, hour, minute int
		expectedChart                  *BaziChart
		description                    string
	}{
		{"Asia/Ho_Chi_Minh",
			1977, 7, 12, 23, 30,
			&BaziChart{
				YearPillar: &YearPillar{
					HeavenlyStem:  HeavenlyStem{Name: YinFireName, Value: YinFireValue},
					EarthlyBranch: EarthlyBranch{Name: Snake, Value: SnakeValue},
					Year:          1977,
				},
				MonthPillar: &MonthPillar{
					HeavenlyStem:  HeavenlyStem{Name: YinFireName, Value: YinFireValue},
					EarthlyBranch: EarthlyBranch{Name: Goat, Value: GoatValue},
					Month:         7,
				},
				DayPillar: &DayPillar{
					HeavenlyStem:  HeavenlyStem{Name: YinMetalName, Value: YinMetalValue},
					EarthlyBranch: EarthlyBranch{Name: Goat, Value: GoatValue},
					Day:           12,
				},
				HourPillar: &HourPillar{
					HeavenlyStem:  HeavenlyStem{Name: YangEarthName, Value: YangEarthValue},
					EarthlyBranch: EarthlyBranch{Name: Rat, Value: RatValue},
					Hour:          TimeOfDay{Hour: 23, Minute: 30},
				},
			},
			"case 1"},
		{"Asia/Ho_Chi_Minh",
			1995, 6, 8, 22, 05,
			&BaziChart{
				YearPillar: &YearPillar{
					HeavenlyStem:  HeavenlyStem{Name: YinWoodName, Value: YinWoodValue},
					EarthlyBranch: EarthlyBranch{Name: Pig, Value: PigValue},
					Year:          1995,
				},
				MonthPillar: &MonthPillar{
					HeavenlyStem:  HeavenlyStem{Name: YangWaterName, Value: YangWaterValue},
					EarthlyBranch: EarthlyBranch{Name: Horse, Value: HorseValue},
					Month:         6,
				},
				DayPillar: &DayPillar{
					HeavenlyStem:  HeavenlyStem{Name: YangMetalName, Value: YangMetalValue},
					EarthlyBranch: EarthlyBranch{Name: Horse, Value: HorseValue},
					Day:           8,
				},
				HourPillar: &HourPillar{
					HeavenlyStem:  HeavenlyStem{Name: YangFireName, Value: YangFireValue},
					EarthlyBranch: EarthlyBranch{Name: Rat, Value: RatValue},
					Hour:          TimeOfDay{Hour: 22, Minute: 05},
				},
			},
			"case 2"},
	}

	// Loop through each test case
	for _, tc := range testCases {
		// Create time object
		location, _ := time.LoadLocation(tc.loc)
		dateTime := time.Date(tc.year, time.Month(tc.month), tc.day, tc.hour, tc.minute, 0, 0, location)

		// Get the chart
		baziChart, err := GetFourPillarChart("", dateTime, location)
		if err != nil {
			t.Errorf("Error for %s: %v", tc.description, err) // Include description in error message
			continue                                          // Skip to next case if this one failed
		}

		jsonData, _ := json.Marshal(baziChart)
		expectData, _ := json.Marshal(tc.expectedChart)

		// Compare JSON output
		if !bytes.Equal(jsonData, expectData) {
			t.Errorf("Mismatch for %s:\nExpected: %v\nGot: %v", tc.description, tc.expectedChart, baziChart)
		}
	}
}
