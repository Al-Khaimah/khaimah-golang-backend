package users

import (
	"time"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/google/uuid"
)

type IamAuth struct {
	base.Model
	UserID         uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"user_id"`
	User           *User     `gorm:"foreignKey:UserID" json:"user"`
	Password       string    `gorm:"type:varchar(255)" json:"password"`
	OTP            int       `gorm:"type:int" json:"otp"`
	ExpiresAt      time.Time `gorm:"type:timestamp" json:"expires_at"`
	FailedAttempts int       `gorm:"type:int" json:"failed_attempts"`
	IsActive       bool      `json:"is_active"`
}
