package main

import (
	"ecart/src/database"
	"ecart/src/models"

	"github.com/bxcodec/faker/v4"
)

func main() {
	database.Connect()
	for i := 0; i < 30; i++ {
		p := models.Product{
			Name:        faker.DomainName(),
			Description: faker.Paragraph(),
			Price:       10,
		}

		database.DB.Create(&p)
	}
}
