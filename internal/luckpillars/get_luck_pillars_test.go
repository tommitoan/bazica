package luckpillars

import (
	"encoding/json"
	"fmt"
	"github.com/tommitoan/bazica/model"
	"reflect"
	"testing"
)

func TestGetLuckPillars(t *testing.T) {
	type args struct {
		fourPillars *model.FourPillars
		gender      int
		prefixPath  []string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.LuckPillars
		wantErr bool
	}{
		{"happy case", args{
			fourPillars: &model.FourPillars{
				YearPillar: &model.YearPillar{
					HeavenlyStem:  model.HeavenlyStem{Name: model.YinWoodName, Value: model.YinWoodValue},
					EarthlyBranch: model.EarthlyBranch{Name: model.Pig, Value: model.PigValue},
					Year:          1995,
				},
				MonthPillar: &model.MonthPillar{
					HeavenlyStem:  model.HeavenlyStem{Name: model.YangWaterName, Value: model.YangWaterValue},
					EarthlyBranch: model.EarthlyBranch{Name: model.Horse, Value: model.HorseValue},
					Month:         6,
				},
				DayPillar: &model.DayPillar{
					HeavenlyStem:  model.HeavenlyStem{Name: model.YangMetalName, Value: model.YangMetalValue},
					EarthlyBranch: model.EarthlyBranch{Name: model.Horse, Value: model.HorseValue},
					Day:           8,
				},
				HourPillar: &model.HourPillar{
					HeavenlyStem:  model.HeavenlyStem{Name: model.YinFireName, Value: model.YinFireValue},
					EarthlyBranch: model.EarthlyBranch{Name: model.Pig, Value: model.PigValue},
					Hour:          model.TimeOfDay{Hour: 22, Minute: 05},
				},
			},
			gender:     1,
			prefixPath: nil,
		}, &model.LuckPillars{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLuckPillars(tt.args.fourPillars, tt.args.gender, tt.args.prefixPath...)
			jsonData, _ := json.Marshal(got)
			fmt.Println(string(jsonData))

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLuckPillars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLuckPillars() got = %v, want %v", got, tt.want)
			}
		})
	}
}
