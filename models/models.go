package models

import (
	"fmt"
)

type UpElement int

const (
	YangWood UpElement = iota + 1
	YinWood
	YangFire
	YinFire
	YangEarth
	YinEarth
	YangMetal
	YinMetal
	YangWater
	YinWater
)

func (e UpElement) String() string {
	switch e {
	case YangWood:
		return "YangWood"
	case YinWood:
		return "YinWood"
	case YangFire:
		return "YangFire"
	case YinFire:
		return "YinFire"
	case YangEarth:
		return "YangEarth"
	case YinEarth:
		return "YinEarth"
	case YangMetal:
		return "YangMetal"
	case YinMetal:
		return "YinMetal"
	case YangWater:
		return "YangWater"
	case YinWater:
		return "YinWater"
	default:
		return fmt.Sprintf("Unknown Element (%d)", e)
	}
}

type DownZodiac int

const (
	Rat DownZodiac = iota + 1
	Ox
	Tiger
	Rabbit
	Dragon
	Snake
	Horse
	Goat
	Monkey
	Rooster
	Dog
	Pig
)

func (z DownZodiac) String() string {
	switch z {
	case Rat:
		return "Rat"
	case Ox:
		return "Ox"
	case Tiger:
		return "Tiger"
	case Rabbit:
		return "Rabbit"
	case Dragon:
		return "Dragon"
	case Snake:
		return "Snake"
	case Horse:
		return "Horse"
	case Goat:
		return "Goat"
	case Monkey:
		return "Monkey"
	case Rooster:
		return "Rooster"
	case Dog:
		return "Dog"
	case Pig:
		return "Pig"
	default:
		return fmt.Sprintf("Unknown Zodiac (%d)", z)
	}
}
