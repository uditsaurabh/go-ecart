package models

type Product struct {
	Id          uint `gorm:"primary_key"`
	Name        string
	Description string
	Price       float64
}
