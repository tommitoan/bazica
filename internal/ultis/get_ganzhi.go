package ultis

import "strings"

var baziMap = map[string]string{
	"Yang Wood Rat":     "Gold in the Sea",
	"Yin Wood Ox":       "Gold in the Sea",
	"Yang Fire Tiger":   "Fire in the Furnace",
	"Yin Fire Rabbit":   "Fire in the Furnace",
	"Yang Earth Dragon": "Wood in the Great Forest",
	"Yin Earth Snake":   "Wood in the Great Forest",
	"Yang Metal Horse":  "Earth on the Roadside",
	"Yin Metal Goat":    "Earth on the Roadside",
	"Yang Water Monkey": "Gold on the Point of the Sword",
	"Yin Water Rooster": "Gold on the Point of the Sword",
	"Yang Wood Dog":     "Fire on the Mountain",
	"Yin Wood Pig":      "Fire on the Mountain",
	"Yang Fire Rat":     "Stream Water",
	"Yin Fire Ox":       "Stream Water",
	"Yang Earth Tiger":  "City Wall Earth",
	"Yin Earth Rabbit":  "City Wall Earth",
	"Yang Metal Dragon": "Waxing Gold",
	"Yin Metal Snake":   "Waxing Gold",
	"Yang Water Horse":  "Willow Tree Wood",
	"Yin Water Goat":    "Willow Tree Wood",
	"Yang Wood Monkey":  "Well Spring Water",
	"Yin Wood Rooster":  "Well Spring Water",
	"Yang Fire Dog":     "Earth on the Roof",
	"Yin Fire Pig":      "Earth on the Roof",
	"Yang Earth Rat":    "Thunder Fire",
	"Yin Earth Ox":      "Thunder Fire",
	"Yang Metal Tiger":  "Pine and Cypress Wood",
	"Yin Metal Rabbit":  "Pine and Cypress Wood",
	"Yang Water Dragon": "Flowing River Water",
	"Yin Water Snake":   "Flowing River Water",
	"Yang Wood Horse":   "Gold in the Sand",
	"Yin Wood Goat":     "Gold in the Sand",
	"Yang Fire Monkey":  "Fire at the Foot of the Mountain",
	"Yin Fire Rooster":  "Fire at the Foot of the Mountain",
	"Yang Earth Dog":    "Wood in the Flatland",
	"Yin Earth Pig":     "Wood in the Flatland",
	"Yang Metal Rat":    "Earth on the Wall",
	"Yin Metal Ox":      "Earth on the Wall",
	"Yang Water Tiger":  "Gold Foil",
	"Yin Water Rabbit":  "Gold Foil",
	"Yang Wood Dragon":  "Lamp Fire",
	"Yin Wood Snake":    "Lamp Fire",
	"Yang Fire Horse":   "Milky Way Water",
	"Yin Fire Goat":     "Milky Way Water",
	"Yang Earth Monkey": "Posthouse Earth",
	"Yin Earth Rooster": "Posthouse Earth",
	"Yang Metal Dog":    "Ornamental Gold",
	"Yin Metal Pig":     "Ornamental Gold",
	"Yang Water Rat":    "Mulberry Wood",
	"Yin Water Ox":      "Mulberry Wood",
	"Yang Wood Tiger":   "Big Stream Water",
	"Yin Wood Rabbit":   "Big Stream Water",
	"Yang Fire Dragon":  "Earth in the Sand",
	"Yin Fire Snake":    "Earth in the Sand",
	"Yang Earth Horse":  "Heavenly Fire",
	"Yin Earth Goat":    "Heavenly Fire",
	"Yang Metal Monkey": "Pomegranate Wood",
	"Yin Metal Rooster": "Pomegranate Wood",
	"Yang Water Dog":    "Ocean Water",
	"Yin Water Pig":     "Ocean Water",
}

// Function to get the corresponding element
func GetGanzhi(sign string) (string, string, int) {
	if element, exists := baziMap[sign]; exists {
		mainElement := ""
		valueElement := 0
		switch {
		case strings.Contains(element, "Fire"):
			mainElement = "Fire"
			valueElement = 2
		case strings.Contains(element, "Wood"):
			mainElement = "Wood"
			valueElement = 3
		case strings.Contains(element, "Water"):
			mainElement = "Water"
			valueElement = 1
		case strings.Contains(element, "Earth"):
			mainElement = "Earth"
			valueElement = 5
		case strings.Contains(element, "Gold"):
			mainElement = "Metal"
			valueElement = 4
		default:
			mainElement = ""
			valueElement = 0
		}
		return element, mainElement, valueElement
	}
	return "", "", 0
}
