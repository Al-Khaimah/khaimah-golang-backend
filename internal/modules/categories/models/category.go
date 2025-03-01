package categories

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
)

type Category struct {
	base.Model
	Name        string `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}
