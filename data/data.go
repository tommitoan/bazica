// Package data embeds the calendar tables used by bazica so the library works
// without any files on disk.
package data

import _ "embed"

// SolarTerm is the JSON table of the 24 solar terms for every year.
//
//go:embed solar-term.json
var SolarTerm []byte

// LunarNewYear is the JSON table of Lunar New Year dates for every year.
//
//go:embed lunar-new-year.json
var LunarNewYear []byte
