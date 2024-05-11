[![License](https://img.shields.io/badge/license-MIT-cyan)](https://github.com/tommitoan/bazica/blob/master/LICENSE)
# Bazica (Ba-zi Calculator) - (In Progress 🚧)
Convert Solar Calendar to Bazi Chart (Chinese astrology) with the year, month, day and hour of birth information (in Go)

# Getting started with bazica

## Prerequisites

- **[Go](https://go.dev/)**: any one of the **three latest major** [releases](https://go.dev/doc/devel/release) (we test it with these).

## Getting Bazica (Ba-zi Calculator)

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/tommitoan/bazica"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `gin` package:

```sh
$ go get -u github.com/tommitoan/bazica
```

## Get Data (Important) ⚠️⚠️

This module provides data related to solar terms, zodiac signs, and the 60-year era cycle.

For your convenience, the data is also summarized in JSON files and located in the `data` folder.

To use this module:

Clone or copy this repository to your project directory.
Access the JSON files located in the data folder.
Please note: You'll need to manually copy the relevant JSON files from the data folder into your project to utilize the data within your application.

### 
https://www.geomancy.net/forums/topic/10229-understand-the-chinese-lunar-and-xia-calendar-in-ba-zi-four-pillars-used-by-various-masters-and-why-not-to-totally-depend-on-just-the-xia-hsia-seasonal-solar-calendar-alone/

For a given date and time of birth, the Hsia Calendar is used to obtain the Heavenly Stems and Earthly Branches with which to construct the Four Pillars of Destiny (四柱) of a person. The Four Pillars of Destiny are also known as the Ba Zi (八字), literally the Eight Characters of a person's birth. In ancient China, one use of this form of divination was to select a suitable marriage partner.

The Four Pillars of Destiny are the Year Pillar, the Month Pillar, the Day Pillar and the Hour Pillar. Each Pillar is made up of two components, namely the Heavenly Stem and the Earthly Branch, both of which are read from the Hsia Calendar.

The Hsia Calendar is also used in various types of Feng Shui, or Chinese Geomancy, such as the Flying Star Feng Shui.

A ba zi chart done with mixed Lunar & Seasonal (Xia) Calendar will look like this:- 

Year: Ji-Chou (Ox) - (Lunar Calendar Year)
Month: Yi-Chou - (Seasonal-Xia Calendar Month)
Day: Xin-Wei - (Lunar Calendar Day) 