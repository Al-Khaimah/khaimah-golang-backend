package users

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
)

type User struct {
	base.Model
	FirstName string `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string `gorm:"type:varchar(100)" json:"last_name"`
	Email     string `gorm:"type:varchar(255);uniqueIndex" json:"email"`

	Categories []categories.Category `gorm:"many2many:user_categories" json:"categories"`
	Bookmarks  []podcasts.Podcast    `gorm:"many2many:user_bookmarks" json:"bookmarks"`
	Auth       IamAuth               `gorm:"foreignKey:UserID" json:"auth"`
}
