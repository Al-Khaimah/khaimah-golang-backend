package users

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        string `json:"id" gorm:"primary_key"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
