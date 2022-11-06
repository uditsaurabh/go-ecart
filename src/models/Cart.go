package models

type Cart struct {
	Id       uint      `gorm:"primaryKey"`
	Products []Product `gorm:"many2many:cart_product;ForeignKey:Id"`
}
