package analysis

import "testing"

// Stems in value order: Jia Yi Bing Ding Wu Ji Geng Xin Ren Gui.
var tenGodRows = map[int][10]string{
	1: {"friend", "rob_wealth", "eating_god", "hurting_officer", "indirect_wealth", "direct_wealth",
		"seven_killings", "direct_officer", "indirect_resource", "direct_resource"},
	2: {"rob_wealth", "friend", "hurting_officer", "eating_god", "direct_wealth", "indirect_wealth",
		"direct_officer", "seven_killings", "direct_resource", "indirect_resource"},
	9: {"eating_god", "hurting_officer", "indirect_wealth", "direct_wealth", "seven_killings", "direct_officer",
		"indirect_resource", "direct_resource", "friend", "rob_wealth"},
}

func TestTenGodAgainstClassicalRows(t *testing.T) {
	for day, row := range tenGodRows {
		for i, want := range row {
			got := tenGod(day, i+1)
			if got.Code != "ten_god."+want {
				t.Errorf("Day Master %d, stem %d: got %s, want ten_god.%s", day, i+1, got.Code, want)
			}
		}
	}
}

func TestTenGodSameStemIsAlwaysFriend(t *testing.T) {
	for s := 1; s <= 10; s++ {
		if got := tenGod(s, s).Code; got != "ten_god.friend" {
			t.Errorf("stem %d against itself = %s, want friend", s, got)
		}
	}
}

func TestTenGodTermsAreBilingual(t *testing.T) {
	got := tenGod(7, 8) // Geng Day Master, Xin stem
	if got.Code != "ten_god.rob_wealth" || got.EN != "Rob Wealth" || got.VI != "Kiếp Tài" {
		t.Errorf("tenGod(7, 8) = %+v", got)
	}
}
