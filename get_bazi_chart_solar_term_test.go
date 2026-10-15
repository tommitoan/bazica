package bazica

import (
	"testing"
	"time"

	"github.com/tommitoan/bazica/v2/model"
)

func TestBirthSolarTermInChart(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		t.Skipf("timezone database unavailable: %v", err)
	}
	cases := []struct {
		name, wantCode, wantVI string
		birth                  time.Time
	}{
		{"A", "solar_term.grain_in_ear", "Mang Chủng", time.Date(1995, 6, 8, 11, 15, 0, 0, loc)},
		{"B", "solar_term.start_of_summer", "Lập Hạ", time.Date(1995, 5, 19, 5, 1, 0, 0, loc)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			chart, err := GetBaziChart(tc.birth, loc, model.GenderFemale)
			if err != nil {
				t.Fatal(err)
			}
			term := chart.Analysis.BirthSolarTerm
			if term == nil {
				t.Fatal("birth_solar_term is nil")
			}
			if term.Term.Code != tc.wantCode || term.Term.VI != tc.wantVI || term.Term.EN == "" {
				t.Errorf("term = %+v, want %s/%s", term.Term, tc.wantCode, tc.wantVI)
			}
			if tc.birth.Before(term.StartedAt) || !tc.birth.Before(term.NextAt) {
				t.Errorf("birth %s outside [%s, %s)", tc.birth, term.StartedAt, term.NextAt)
			}
		})
	}
}

func TestBirthSolarTermUsesChartLocation(t *testing.T) {
	loc := time.FixedZone("Vietnam", 7*60*60)
	birth := time.Date(1995, 6, 8, 4, 15, 0, 0, time.UTC)
	utc, err := GetBaziChart(birth, time.UTC, model.GenderMale)
	if err != nil {
		t.Fatal(err)
	}
	local, err := GetBaziChart(birth, loc, model.GenderMale)
	if err != nil {
		t.Fatal(err)
	}
	a, b := utc.Analysis.BirthSolarTerm, local.Analysis.BirthSolarTerm
	if a.Term.Code != b.Term.Code || !a.StartedAt.Equal(b.StartedAt) || !a.NextAt.Equal(b.NextAt) {
		t.Errorf("different instant: UTC %+v, local %+v", a, b)
	}
	if _, offset := b.StartedAt.Zone(); offset != 7*60*60 {
		t.Errorf("local term offset = %d, want +7h", offset)
	}
}
