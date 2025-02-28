package notifications

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	"github.com/google/uuid"
)

type Notification struct {
	base.Model
	UserID      uuid.UUID   `gorm:"type:uuid;index" json:"user_id"`
	User        *users.User `gorm:"foreignKey:UserID" json:"user"`
	Title       string      `gorm:"type:varchar(255)" json:"title"`
	Description string      `gorm:"type:text" json:"description"`
	IsRead      bool        `gorm:"default:false" json:"is_read"`
	Type        string      `gorm:"type:varchar(100)" json:"type"`
}
