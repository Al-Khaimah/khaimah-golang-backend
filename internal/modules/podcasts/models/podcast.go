package podcasts

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/google/uuid"
)

type Podcast struct {
	base.Model
	Title       string    `gorm:"type:varchar(255);index" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	ContentURL  string    `gorm:"type:text" json:"content_url"`
	LikesCount  int       `gorm:"default:0" json:"likes_count"`
	CategoryID  uuid.UUID `gorm:"type:uuid;index" json:"category_id"`
}

type UserPodcast struct {
	base.Model
	UserID         uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	PodcastID      uuid.UUID `gorm:"type:uuid;index" json:"podcast_id"`
	CategoryID     uuid.UUID `gorm:"type:uuid;index" json:"category_id"`
	ResumePosition int       `gorm:"default:0" json:"resume_position"`
	IsCompleted    bool      `gorm:"default:false" json:"is_completed"`
}
