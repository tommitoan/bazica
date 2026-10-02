[![CI](https://github.com/tommitoan/bazica/actions/workflows/ci.yml/badge.svg)](https://github.com/tommitoan/bazica/actions/workflows/ci.yml)
[![GitHub release](https://img.shields.io/github/tag/tommitoan/bazica.svg?label=latest)](https://github.com/tommitoan/bazica/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/tommitoan/bazica.svg)](https://pkg.go.dev/github.com/tommitoan/bazica)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommitoan/bazica)](https://goreportcard.com/report/github.com/tommitoan/bazica)
[![License](https://img.shields.io/badge/license-MIT-cyan)](https://github.com/tommitoan/bazica/blob/main/LICENSE)

<p align="center">
  <img style="float: right;" src="./Images/bazica-gopher.png" width=400>
</p>

# Bazica (Ba-zi Chart Calculator) 
Convert Solar Calendar to Bazi Chart (Chinese astrology) with the year, month, day and hour of birth information (in Go)

<div align="center">
  <img alt="Demo" src="./Images/bazica-screen.png" />
</div>

<h3 align="center">
  <a href="https://bazica.onrender.com/" target="_blank">Live Demo</a>
</h3>

## Getting started with bazica 
### Prerequisites

- **[Go](https://go.dev/)**: 1.21 or newer. CI runs the tests on Go 1.21 and on the two most recent [releases](https://go.dev/doc/devel/release).

### Importing Bazica

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/tommitoan/bazica"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `bazica` package:

```sh
go get -u github.com/tommitoan/bazica@latest
```

## Data

The solar term and Lunar New Year tables are embedded in the module (`go:embed`), so nothing has to be copied into your project and the library works from any working directory, including a vendored build or a container image with only the binary.

The same tables are also published as JSON in the `data` folder if you want to use them elsewhere:

- `data/solar-term.json`: the 24 solar terms for each year, in UTC.
- `data/lunar-new-year.json`: the Lunar New Year date for each year.

## How to use

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/tommitoan/bazica"
	"time"
)

func main() {
	// Calculate current ba-zi chart
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Now()
	gender := 0 // 0 = female & 1 = male
	
	/* Example: 
	Time to calculate: 1990-12-31 6:30 - Timezone: HoChiMinh / Vietnam - Gender: Male
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	now := time.Date(1990, time.Month(12), 31, 6, 30, 0, 0, loc)
	gender := 1
	*/

	chart, err := bazica.GetBaziChart(now, loc, gender)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonData, _ := json.Marshal(chart)
	fmt.Println(string(jsonData))
}
```
## Note
### Data Input Limitations:
Due to the specific calculations and algorithms used in this package, it is currently designed to handle date inputs ranging from January 1, 1900, to December 31, 2099. Dates outside this range return `model.ErrDateOutOfRange`.   

### Conventions
- **Day**: the day changes at 23:00 (the Rat hour), so a birth at 23:30 takes the pillar of the next day.
- **Year**: the year changes at the Lunar New Year (`data/lunar-new-year.json`).
- **Month**: the month changes at the "initial" solar terms (jie), such as Start of Spring and Awakening of Insects.
- **Time zone**: `dateTime` is read as wall-clock time in `loc`, so pass the birth place's zone. A nil `loc` keeps the zone `dateTime` already carries.
- **Gender**: `0` is female and `1` is male (`model.GenderFemale`, `model.GenderMale`); it sets the direction of the luck pillars. Other values return `model.ErrInvalidGender`.

If you require calculations for dates outside this range, please consider alternative libraries or solutions.

###  Ho Chi Minh City 1975
Ho Chi Minh City (formerly Saigon) has gone through time zone changes throughout history. Notably:  
Before 1975: South Vietnam (including Saigon) used UTC+8.  
After 1975: The unified Vietnam adopted UTC+7.

## References

This project drew inspiration and information from the following sources:

* **[Thời Gian](https://www.thoigian.com.vn/)** - A comprehensive resource for understanding the Vietnamese calendar system.
* **[Chinese Fortune Calendar](https://www.chinesefortunecalendar.com/)** - Provided insights into the Chinese calendar, calculations, and cultural significance.


### Document
https://www.geomancy.net/forums/topic/10229-understand-the-chinese-lunar-and-xia-calendar-in-ba-zi-four-pillars-used-by-various-masters-and-why-not-to-totally-depend-on-just-the-xia-hsia-seasonal-solar-calendar-alone/

The article above discusses the Lunar and Hsia (seasonal, solar term) calendars used to read the Four Pillars. Bazica uses a mix of both: the year follows the lunar calendar and the month follows the solar terms. An example chart:

```
Date time: 2024-07-05 22:00 UTC+7 (Vietnam)  

Year: Yang Wood - Dragon (Lunar Calendar Year)  
Month: Yang Metal - Horse (Solar-Term Calendar Month)  
Day: Yang Metal - Horse (Lunar Calendar Day)   
Hour: Yin Wood - Pig (Zodiac Hour)
```


