package bazica

import (
	"embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tommitoan/bazica/svc/datasvc"
	toerr "github.com/tommitoan/bazica/svc/toerr"
	"net/http"
	"os"
	"testing"
)

func TestBazicaSvc_GetSolarTermsByYear(t *testing.T) {
	e, _ := os.Getwd()
	fmt.Println(e)
	type fields struct {
		DataSvc *datasvc.DataSvc
		Solar   *datasvc.SolarTerms
	}
	type args struct {
		prefix string
		year   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:   "happy case: 1990",
			fields: fields{},
			args: args{
				prefix: "svc/datasvc/solar-terms-data/",
				year:   "1990",
			},
			want:    "1990-01-05 14:33:14.870+00:00",
			wantErr: nil,
		},
		{
			name:   "error patch",
			fields: fields{},
			args: args{
				prefix: "wrong/path/",
				year:   "1990",
			},
			want:    "1990-01-05 14:33:14.870+00:00",
			wantErr: toerr.NewValidationError(http.StatusInternalServerError, "Error when open file"),
		},
		{
			name:   "no support year",
			fields: fields{},
			args: args{
				prefix: "svc/datasvc/solar-terms-data/",
				year:   "2150",
			},
			want:    "",
			wantErr: toerr.NewValidationError(http.StatusBadRequest, "Year not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bzc := &BazicaSvc{}
			res, gotErr := bzc.GetSolarTermsByYear(tt.args.prefix, tt.args.year)
			if res != nil {
				assert.Equalf(t, tt.want, res.MinorCold, "case(%s)", tt.name)
			}
			if assert.Equalf(t, tt.wantErr, gotErr, "case(%s)", tt.name) {
				return
			}

		})
	}
}

func TestBazicaSvc_GetAllSolarTerms(t *testing.T) {
	type fields struct {
		DataSvc *datasvc.DataSvc
		Solar   *datasvc.SolarTerms
		dataDir embed.FS
	}
	tests := []struct {
		name    string
		fields  fields
		want    [200]datasvc.SolarTerms
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "happy case",
			fields:  fields{},
			want:    [200]datasvc.SolarTerms{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bzc := &BazicaSvc{
				DataSvc: tt.fields.DataSvc,
				Solar:   tt.fields.Solar,
				dataDir: tt.fields.dataDir,
			}
			got, err := bzc.GetAllSolarTerms()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllSolarTerms()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllSolarTerms()")
		})
	}
}
