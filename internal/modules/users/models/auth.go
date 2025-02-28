package users

import (
	"time"

	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	ID        string    `json:"id" gorm:"primary_key"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;references:ID"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID;"`
	Password  string    `json:"password"`
	IsActive  bool      `json:"is_active"`
	LastLogin time.Time `json:"last_login"`
}
