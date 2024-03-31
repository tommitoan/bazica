package main

import (
	"bazica/bazica"
	"bazica/models"
)

func main() {
	bazica.CalTest()
	println(models.Tiger)
	println(models.Tiger.String())

}
