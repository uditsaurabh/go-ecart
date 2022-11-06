package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        uint `json:"-"`
	FirstName string
	LastName  string
	Email     string
	Password  string `json:"-"`
	IsAdmin   bool
	Cart      Cart `gorm:"ForeignKey:Id"`
}

func (user *User) SetPassword(password string) error {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	} else {
		user.Password = string(bs)
	}
	return nil
}
func MatchPassword(password, password_confirm interface{}) bool {
	return password == password_confirm
}
