package database

import (
	"ecart/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:root@tcp(db:3306)/ekart"), &gorm.Config{})
	if err != nil {
		panic("could not connect to db")
	}
}

func Automigrate() {

	DB.AutoMigrate(models.User{}, models.Cart{}, models.Product{})
}
