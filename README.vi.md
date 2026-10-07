[![CI](https://github.com/tommitoan/bazica/actions/workflows/ci.yml/badge.svg)](https://github.com/tommitoan/bazica/actions/workflows/ci.yml)
[![GitHub release](https://img.shields.io/github/tag/tommitoan/bazica.svg?label=latest)](https://github.com/tommitoan/bazica/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/tommitoan/bazica/v2.svg)](https://pkg.go.dev/github.com/tommitoan/bazica/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommitoan/bazica/v2)](https://goreportcard.com/report/github.com/tommitoan/bazica/v2)
[![License](https://img.shields.io/badge/license-MIT-cyan)](https://github.com/tommitoan/bazica/blob/master/LICENSE)

<p align="center">
  <img src="./Images/bazica-gopher.png" width=400 alt="Bazica gopher">
</p>

# Bazica (Lập lá số Bát Tự)

[English](README.md) · **Tiếng Việt**

Bazica là thư viện Go đổi giờ sinh theo dương lịch thành lá số Bát Tự (Tứ Trụ, tử vi phương Đông): trụ năm, tháng, ngày, giờ, các đại vận, các lưu niên, cùng một lớp dữ kiện suy ra (Thập Thần, Trường Sinh, Tàng Can, Nạp Âm, các sao và nhiều hơn nữa). Mọi nhãn đều có cả tiếng Anh và tiếng Việt.

<div align="center">
  <img alt="bazica-web: nhập giờ sinh, xem các trụ, đại vận và một thẻ cho từng năm" src="./Images/bazica-web-demo.gif" width="820" />
  <br>
  <sub>Ứng dụng web đi kèm, <a href="https://github.com/tommitoan/bazica-web">bazica-web</a>, xây trên thư viện này.</sub>
</div>

<h3 align="center">
  <a href="https://bazi.tommitoan.com/" target="_blank">Dùng thử trực tuyến</a>
</h3>

## Điểm nổi bật

- **Bốn trụ**: trụ năm đổi đúng thời điểm Lập Xuân, trụ tháng đổi theo các tiết khí, cho các năm 1700-2399.
- **Đại vận** với chiều thuận/nghịch, thời điểm khởi vận, các năm và tuổi ước tính (tuổi mụ theo Bát Tự).
- **Lưu niên** (`GetAnnualPillars`) kèm Thập Thần, Trường Sinh, đại vận chứa năm đó và cờ xung.
- **Phân tích** trên mọi lá số: Thập Thần, mười hai Trường Sinh, Tàng Can, Nạp Âm, Không Vong, thiên khắc địa xung, 58 sao đã kiểm chứng, Thai Nguyên, Thai Tức, Mệnh Cung và số lượng từng hành.
- **Thuật ngữ song ngữ với mã cố định**: mỗi khái niệm là `{code, en, vi}`, giao diện đổi ngôn ngữ mà không phải tính lại lá số.
- **Không cần kèm file dữ liệu**: bảng tiết khí và bảng Tết Nguyên Đán được nhúng sẵn trong module. Hai bảng này do `tools/gencal` sinh ra từ lịch thiên văn JPL DE440.
- **Kiểm thử kỹ**: kết quả được đối chiếu với một bản cài đặt khác và với các lá số tham chiếu đã ghi lại (xem "Độ chính xác và cách kiểm tra").

## Bắt đầu

### Yêu cầu

- **[Go](https://go.dev/)**: 1.21 trở lên. CI chạy kiểm thử trên Go 1.21 và hai [bản phát hành](https://go.dev/doc/devel/release) mới nhất.

### Cài đặt

```sh
go get github.com/tommitoan/bazica/v2@latest
```

rồi import:

```go
import "github.com/tommitoan/bazica/v2"
```

Đường dẫn module kết thúc bằng `/v2` (Go bắt buộc có hậu tố này từ v2 trở đi). Import `github.com/tommitoan/bazica` không có `/v2` vẫn dừng ở v1.5.0; xem mục "Phiên bản 2" bên dưới để biết cái gì đã đổi.

### Lá số đầu tiên

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tommitoan/bazica/v2"
	"github.com/tommitoan/bazica/v2/model"
)

func main() {
	// 1990-12-31 06:30, TP. Hồ Chí Minh, nam.
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatal(err)
	}
	birth := time.Date(1990, time.December, 31, 6, 30, 0, 0, loc)

	chart, err := bazica.GetBaziChart(birth, loc, model.GenderMale)
	if err != nil {
		log.Fatal(err) // ví dụ model.ErrDateOutOfRange
	}

	p := chart.FourPillar
	fmt.Println("year: ", p.YearPillar.HeavenlyStem.Name, p.YearPillar.EarthlyBranch.Name)
	fmt.Println("month:", p.MonthPillar.HeavenlyStem.Name, p.MonthPillar.EarthlyBranch.Name)
	fmt.Println("day:  ", p.DayPillar.HeavenlyStem.Name, p.DayPillar.EarthlyBranch.Name)
	fmt.Println("hour: ", p.HourPillar.HeavenlyStem.Name, p.HourPillar.EarthlyBranch.Name)
	fmt.Println("day master:", chart.Analysis.DayMaster.HeavenlyStem.Name, "/", chart.Analysis.DayMaster.Element.VI)

	for _, lp := range chart.LuckPillars.LuckPillars[:3] {
		fmt.Printf("luck pillar %d: %s %s, %d-%d\n", lp.Number, lp.HeavenlyStem.Name, lp.EarthlyBranch.Name, lp.YearStart, lp.YearEnd)
	}

	annual, err := bazica.GetAnnualPillars(chart, 2026, 2)
	if err != nil {
		log.Fatal(err)
	}
	for _, y := range annual.AnnualPillars {
		fmt.Println(y.Year, y.NominalAge, y.HeavenlyStem.Name, y.EarthlyBranch.Name, y.TenGod.EN)
	}
}
```

Kết quả (tên Thiên Can/Địa Chi và Thập Thần ở đây là nhãn tiếng Anh; nhãn tiếng Việt nằm trong các trường `{code, en, vi}` của kết quả phân tích):

```
year:  Yang Metal Horse
month: Yang Earth Rat
day:   Yang Metal Horse
hour:  Yin Earth Rabbit
day master: Yang Metal / Kim
luck pillar 0: Yang Earth Rat, 1990-1991
luck pillar 1: Yin Earth Ox, 1992-2001
luck pillar 2: Yang Metal Tiger, 2002-2011
2026 37 Yang Fire Horse Seven Killings
2027 38 Yin Fire Goat Direct Officer
```

Tức là Canh Ngọ, Mậu Tý, Canh Ngọ, Kỷ Mão, Nhật Chủ Canh (Kim, Dương).

`GetBaziChart` trả về `*model.BaziChart` có thể chuyển thành JSON (`json.Marshal(chart)`), nên dùng làm phần lõi của một API cũng được. Đại vận số `0` là trụ tháng sinh, áp dụng cho đến khi đại vận đầu tiên bắt đầu. Các ví dụ chạy được, có kiểm tra kết quả, nằm trong `example_test.go` và trên [pkg.go.dev](https://pkg.go.dev/github.com/tommitoan/bazica/v2).

### Phạm vi ngày sinh và dữ liệu

Hỗ trợ ngày sinh từ 01/01/1700 đến 31/12/2399. `bazica.SupportedYears()` trả về năm đầu và năm cuối, lấy từ các bảng lịch nhúng sẵn; ngày ngoài phạm vi trả về `model.ErrDateOutOfRange`.

Bảng tiết khí và bảng Tết Nguyên Đán được nhúng trong module (`go:embed`), nên không phải chép gì vào dự án và thư viện chạy được từ bất kỳ thư mục nào, kể cả bản build vendor hay container chỉ có file thực thi. Hai bảng này cũng được công bố dạng JSON trong thư mục `data` nếu bạn muốn dùng ở nơi khác:

- `data/solar-term.json`: 24 tiết khí của từng năm, tính theo UTC.
- `data/lunar-new-year.json`: ngày Tết Nguyên Đán của từng năm. Bảng này chỉ cấp dữ liệu cho khối `year_reference` (chỉ để hiển thị); không có trụ nào được tính từ nó.

## Phiên bản 2

`v2.0.0` đổi đường dẫn module thành `github.com/tommitoan/bazica/v2` (Go bắt buộc có hậu tố này từ v2); không hàm hay kiểu nào bị đổi tên hay xóa, nên nâng cấp chỉ là đổi đường dẫn import rồi chạy `go get github.com/tommitoan/bazica/v2@latest`. v1.5.0 vẫn còn đó cho ai muốn giữ cách tính cũ.

**Kết quả thay đổi thế nào.** Trụ năm giờ đổi đúng thời điểm Lập Xuân thay vì Tết Nguyên Đán, giống các công cụ Bát Tự theo tiết khí. Lá số của người sinh ngoài hai khoảng dưới đây không đổi (một bài kiểm tra so sánh 658 lá số với v1.5.0). Với người sinh giữa Tết Nguyên Đán và Lập Xuân, hoặc giữa Lập Xuân và Tết Nguyên Đán (khoảng 2 % số người, trung bình 7,4 ngày mỗi năm, nhiều nhất 16 ngày), các phần sau thay đổi:

- trụ năm và mọi thứ suy ra từ nó (Thập Thần, Tàng Can, Trường Sinh và Nạp Âm của trụ năm);
- các sao: một sao đổi khi quy tắc của nó xuất phát từ trụ năm hoặc khi nó rơi vào trụ năm hay trụ tháng, tức là phần lớn bảng sao. Khi kiểm tra mọi ngày như vậy trong 1901-2099, cả hai giới, **mọi** lá số đều có ít nhất một sao khác và 39 tên sao có liên quan, vì vậy hãy xem các sao của những người sinh này là mới;
- Thiên Can của tháng (Địa Chi của tháng vẫn theo tiết khí như trước);
- chiều và thứ tự của đại vận, vì âm dương của Thiên Can năm quyết định chiều đi, cùng dòng đầu tiên của bảng lưu niên với tuổi ước tính.

Trụ ngày và trụ giờ không bao giờ đổi. Kết quả được đối chiếu với một bản cài đặt thứ hai (lunar-javascript) cho mọi ngày như vậy trong 1901-2099, cả hai giới: bốn trụ, đại vận đầu tiên và chiều đại vận (`testdata/calendar/lunarjs-window-1901-2099.json`).

**Phần được thêm.** `analysis.year_reference` báo, cho mọi lá số và chỉ để hiển thị, thời điểm Lập Xuân của năm dương lịch của người sinh (theo múi giờ nơi sinh), năm âm lịch cùng Thiên Can, Địa Chi, con giáp và ngày Tết Nguyên Đán, và `differs`, bằng true khi năm âm lịch đó khác trụ năm. Không có gì được tính từ khối này; có một bài kiểm tra bảo đảm điều đó.

Năm âm lịch theo bảng Tết Nguyên Đán, mỗi năm một ngày và theo lịch Trung Quốc (UTC+8), bất kể `loc` bạn truyền vào. Lịch Việt Nam (UTC+7) đặt Tết vào một ngày khác trong vài năm: trong 1900-2099 là 1903, 1935, 1965, 1968, 1969, 1985, 2007, 2030 và 2053 (`tools/gencal/README.md` có bảng so sánh; các năm sau đó chưa được kiểm tra). Vào đúng ngày đó của những năm này, `lunar_year` và `differs` có thể khác lịch vạn niên Việt Nam. Không trụ nào bị ảnh hưởng, vì trụ năm theo Lập Xuân.

## Phân tích

Từ v1.4.0 mỗi lá số còn mang các dữ kiện suy ra dưới khóa `analysis`. Các trường có từ v1.3.0 không đổi, và mọi khóa thêm vào đều tên `analysis`, nên các client cũ vẫn chạy bình thường.

| Ở đâu | Trường | Ý nghĩa |
|---|---|---|
| mỗi trụ của lá số | `ten_god` | Thập Thần của Thiên Can so với Nhật Chủ (`null` ở trụ ngày) |
| | `stem_element`, `branch_element`, `nayin` | ngũ hành và Nạp Âm |
| | `stem_stage_at_own_branch`, `stem_stage_at_month_branch`, `day_master_stage` | mười hai Trường Sinh |
| | `hidden_stems` | các Can ẩn trong Địa Chi, khí chính trước, mỗi Can kèm Thập Thần và Trường Sinh |
| | `is_void`, `heaven_earth_clash` | cờ Không Vong và cờ xung |
| | `stars` | 58 sao đã kiểm chứng, theo thứ tự cố định (danh sách rỗng khi không có); mỗi sao có `nature` (xem bên dưới) |
| mỗi đại vận | `age_start`, `age_end` | khoảng tuổi ước tính |
| | `ten_god`, `nayin`, `stem_stage_at_own_branch`, `hidden_stems`, `heaven_earth_clash` | như trên |
| cả lá số | `day_master`, `void_branches` | Nhật Chủ và hai Địa Chi Không Vong |
| | `year_reference` | chỉ để hiển thị: thời điểm Lập Xuân của năm sinh, năm âm lịch (Can, Chi, con giáp, ngày Tết) và `differs` (xem "Phiên bản 2") |
| | `thai_nguyen`, `thai_tuc`, `life_palace` | Thai Nguyên, Thai Tức và Mệnh Cung cùng Nạp Âm của chúng |
| | `element_counts` | số lượng Thiên Can, Địa Chi và Can ẩn theo từng hành, chưa gán trọng số |
| | `day_master_strength`, `useful_god` | để dành, luôn là `null` |

Mỗi khái niệm có tên (Thập Thần, Trường Sinh, hành, âm dương, Nạp Âm, sao) là một đối tượng có `code` cố định cùng nhãn tiếng Anh và tiếng Việt:

```json
{"code": "ten_god.hurting_officer", "en": "Hurting Officer", "vi": "Thương Quan"}
```

Hãy rẽ nhánh theo `code`, đừng theo nhãn. Giao diện có thể hiện ngôn ngữ nào cũng được mà không cần tính lại lá số.

`GetAnnualPillars(chart, fromYear, count)` trả về trụ của từng năm kèm Thập Thần, Trường Sinh, tuổi ước tính, đại vận chứa năm đó và cờ xung:

```go
chart, _ := bazica.GetBaziChart(birth, loc, model.GenderMale)
annual, err := bazica.GetAnnualPillars(chart, 2026, 10)
if err != nil { // model.ErrInvalidYearRange
	return err
}
for _, y := range annual.AnnualPillars {
	fmt.Println(y.Year, y.NominalAge, y.HeavenlyStem.Name, y.EarthlyBranch.Name, y.TenGod.EN)
}
```

`fromYear` không được trước năm của trụ năm hay năm sinh đầu tiên được hỗ trợ (xem `SupportedYears`), và các năm không được vượt 9999. Bảng lưu niên không bị giới hạn bởi phạm vi ngày sinh: trụ mỗi năm là chu kỳ sáu mươi năm, nên nó tiếp tục sau 2399 (số đại vận là `null` khi năm đã qua đại vận thứ mười hai).

### Quy ước phân tích

- **Tuổi ước tính (tuổi mụ theo Bát Tự)**: năm của trụ năm là tuổi 1. Với người sinh trước Lập Xuân, trụ năm thuộc về năm trước, nên tuổi đếm từ đó. Đây là cách đếm theo Bát Tự, không phải cách đếm thường ngày từ Tết Nguyên Đán.
- **Các sao**: chỉ báo những sao có quy tắc đã được kiểm chứng với các lá số tham chiếu (58 sao, kiểm trên 298 lá số đã ghi lại và 10.000 ngày sinh ngẫu nhiên so với một bản cài đặt độc lập). Các sao khác không bao giờ được đoán. Vài quy tắc khác bảng cổ điển vì theo trang tham chiếu, ví dụ Hồng Diễm (`hong_diem`) xét cả Can ngày lẫn Can năm, và Đại Hao (`dai_hao`) phụ thuộc giới tính.
- **Tính chất của sao (`nature`)**: mỗi thuật ngữ sao có một khóa `nature` tùy chọn, chỉ có ở sao: `auspicious` (cát), `inauspicious` (hung) hoặc `mixed` (tùy cục, tác dụng phụ thuộc phần còn lại của lá số). Đây chỉ là gợi ý để hiển thị. Các trường phái bất đồng về điều gì là tốt và nhiều sao có hai mặt, nên giá trị này không phải lời phán về lá số. Trang tham chiếu chỉ gán sao là tốt hoặc xấu; thư viện theo trang đó, trừ 10 sao được xem là tùy cục và một sao (`am_duong_sat`) được xem là hung trong khi trang ghi là tốt. Các thuật ngữ khác (Thập Thần, Trường Sinh, hành, Nạp Âm) không bao giờ có khóa này. Khóa này chỉ thêm vào: client nào bỏ qua nó vẫn chạy bình thường.
- **Thái Dương (Moon General, `thai_duong`)**: vị tướng đổi ở mỗi trung khí, và cả ngày dương lịch có trung khí rơi vào đã thuộc về vị tướng mới. Trang tham chiếu ghi vài tiết lệch một ngày, nên vào ngày có tiết, kết quả có thể khác trang đó một ngày.
- **Mệnh Cung**: Tý và Sửu được đếm trước Dần khi suy ra Thiên Can, khác với thứ tự tháng cổ điển.
- **Thiên khắc địa xung**: có hướng. Một ô được đánh cờ khi Thiên Can của nó khắc Thiên Can của một trụ khác trong lá số, cùng âm dương, và hai Địa Chi xung nhau. Các Can Thổ cũng tham gia.
- **Cường nhược và Dụng Thần**: chưa tính.

## Quy ước

- **Ngày**: ngày đổi lúc 23:00 (giờ Tý), nên người sinh lúc 23:30 lấy trụ của ngày kế tiếp.
- **Năm**: năm đổi đúng thời điểm Lập Xuân, cùng đồng hồ tiết khí với tháng; phép so sánh dùng thời điểm chính xác, nên quy tắc ngày đổi lúc 23:00 không áp dụng cho nó.
- **Tháng**: tháng đổi ở các tiết "đầu tháng" (tiết, jie), như Lập Xuân và Kinh Trập. Thiên Can của tháng suy từ Thiên Can của năm theo quy tắc Ngũ Hổ Độn.
- **Năm âm lịch (để tham chiếu)**: năm đổi ở Tết Nguyên Đán, báo trong `analysis.year_reference` chỉ để hiển thị. Nó theo ngày dương lịch (âm lịch đổi năm vào nửa đêm, nên quy tắc 23:00 không áp dụng). Với người sinh giữa Tết và Lập Xuân, nó khác trụ năm.
- **Múi giờ**: `dateTime` được đọc là giờ đồng hồ tại `loc`, nên hãy truyền múi giờ của nơi sinh. `loc` nil sẽ giữ múi giờ mà `dateTime` đã mang.
- **Giới tính**: `0` là nữ và `1` là nam (`model.GenderFemale`, `model.GenderMale`); nó quyết định chiều của đại vận. Giá trị khác trả về `model.ErrInvalidGender`.

### Độ chính xác và cách kiểm tra

Hai bảng (`data/solar-term.json` và `data/lunar-new-year.json`) do `tools/gencal` sinh ra từ lịch thiên văn JPL DE440; README của công cụ nói rõ cách làm và độ chính xác. Tết Nguyên Đán theo lịch Trung Quốc (UTC+8, và giờ trung bình Bắc Kinh trước 1929). Độ chính xác của một tiết khí phụ thuộc Delta T, độ chênh giữa thời gian đều và vòng quay của Trái Đất: nó được dựng lại từ quan sát cho đến nay, nên các tiết chính xác đến vài giây từ 1700 đến 2025, còn sau đó là dự báo, nên độ bất định tăng lên đến vài phút vào năm 2399. Một lần đối chiếu với gói độc lập thứ hai (phương pháp khác và Delta T riêng) khớp trong khoảng một phút từ 1700 đến 2000; sau 2100 hai dự báo Delta T lệch nhau 3 đến 5 phút, nên người sinh trong khoảng năm phút quanh một tiết khí là không chắc chắn ở giai đoạn đó, và Tết Nguyên Đán của các năm 2299 và 2333 cũng vậy. Người sinh trong khoảng sai số đó quanh một tiết khí, kể cả Lập Xuân, có thể rơi vào trụ tháng hoặc trụ năm bên cạnh. Tết của 2299 và 2333 cũng bất định như thế, và chúng chỉ ảnh hưởng đến năm âm lịch dùng để hiển thị.

Trước khoảng năm 1900 nhiều nơi dùng giờ trung bình địa phương. `bazica` đọc giờ đồng hồ theo `loc` bạn truyền, nên hãy truyền múi giờ (hoặc độ lệch cố định) đã dùng tại nơi sinh vào lúc đó.

Nếu cần tính cho ngày ngoài phạm vi hỗ trợ, hãy xem xét các thư viện hoặc giải pháp khác.

### TP. Hồ Chí Minh năm 1975

TP. Hồ Chí Minh (trước là Sài Gòn) đã đổi múi giờ nhiều lần trong lịch sử. Đáng chú ý:
Trước 1975: miền Nam Việt Nam (gồm Sài Gòn) dùng UTC+8.
Sau 1975: Việt Nam thống nhất dùng UTC+7.

## Ứng dụng web

[bazica-web](https://github.com/tommitoan/bazica-web) là một máy chủ Go cùng một trang tĩnh xây trên thư viện này. Nó không có logic Bát Tự riêng: nó kiểm tra yêu cầu rồi gọi `GetBaziChart` và `GetAnnualPillars`. Trang hiển thị bốn trụ cùng đầy đủ dữ kiện, sơ đồ Ngũ Hành, bản đồ Thần Sát, các đại vận và một thẻ cho từng năm trong sáu mươi năm, bằng tiếng Việt hoặc tiếng Anh, có liên kết chia sẻ lá số và xuất PDF tùy chọn.

<div align="center">
  <img alt="Các thẻ lưu niên của bazica-web: sáu mươi năm, sáu thẻ mỗi hàng, năm đang xem được tô nổi" src="./Images/bazica-web-yearly-cards.png" width="820" />
</div>

## Phát triển

```sh
go test -race -cover ./...
gofmt -l .
go vet ./...
```

`tools/gencal` sinh lại hai bảng lịch; README của nó giải thích phương pháp và các phép kiểm tra.

## Tài liệu tham khảo

Dự án này lấy cảm hứng và thông tin từ các nguồn sau:

* **[Thời Gian](https://www.thoigian.com.vn/)** - Nguồn đầy đủ để hiểu hệ thống lịch Việt Nam.
* **[Chinese Fortune Calendar](https://www.chinesefortunecalendar.com/)** - Cho hiểu biết về lịch Trung Quốc, cách tính và ý nghĩa văn hóa.
* **[Understand the Chinese Lunar and Xia calendar in Ba-zi](https://www.geomancy.net/forums/topic/10229-understand-the-chinese-lunar-and-xia-calendar-in-ba-zi-four-pillars-used-by-various-masters-and-why-not-to-totally-depend-on-just-the-xia-hsia-seasonal-solar-calendar-alone/)** - Bàn về lịch âm và lịch Hạ (theo mùa, tiết khí) dùng để đọc Tứ Trụ. Từ v2.0.0 Bazica theo lịch tiết khí cho cả năm lẫn tháng; v1 đổi năm ở Tết Nguyên Đán.

## Giấy phép

MIT, xem [LICENSE](LICENSE).
