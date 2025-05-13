package users

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/enums"
)

type User struct {
	base.Model
	FirstName string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string         `gorm:"type:varchar(100)" json:"last_name"`
	UserType  users.UserType `gorm:"type:varchar(20);default:'free'" json:"user_type"`
	Email     string         `gorm:"type:varchar(255);index" json:"email"`

	Categories []categories.Category `gorm:"many2many:user_categories" json:"categories"`
	Bookmarks  []podcasts.Podcast    `gorm:"many2many:user_bookmarks" json:"bookmarks,omitempty"`
	Downloads  []podcasts.Podcast    `gorm:"many2many:user_downloads" json:"downloads,omitempty"`
	Auth       IamAuth               `gorm:"foreignKey:UserID" json:"auth"`
}
