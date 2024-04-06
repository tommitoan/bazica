package datasvc

import (
	"testing"
)

func TestGetSolarTermsByYear(t *testing.T) {

	tests := []struct {
		name        string
		year        string
		expectedErr error
	}{
		// TODO: Add test cases.
		{
			name:        "random",
			year:        "1900",
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := GetSolarTermsByYear("solar-terms-data/", tt.year)
			if gotErr != tt.expectedErr {
				t.Errorf("%s: got = %v, expect = %v", tt.name, gotErr, tt.expectedErr)
				return
			}

		})
	}
}
