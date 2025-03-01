package podcasts

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/google/uuid"
)

type Podcast struct {
	base.Model
	Title                 string    `gorm:"type:varchar(255);index" json:"title"`
	Description           string    `gorm:"type:text" json:"description"`
	AudioURL              string    `gorm:"type:text" json:"audio_url"`
	CoverImageURL         string    `gorm:"type:text" json:"cover_image_url"`
	CoverImageDescription string    `gorm:"type:text" json:"cover_image_description"`
	LikesCount            int       `gorm:"default:0" json:"likes_count"`
	CategoryID            uuid.UUID `gorm:"type:uuid;index" json:"category_id"`
}

type UserPodcast struct {
	base.Model
	UserID         uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	PodcastID      uuid.UUID `gorm:"type:uuid;index" json:"podcast_id"`
	CategoryID     uuid.UUID `gorm:"type:uuid;index" json:"category_id"`
	ResumePosition int       `gorm:"default:0" json:"resume_position"`
	IsCompleted    bool      `gorm:"default:false" json:"is_completed"`
	IsDownloaded   bool      `gorm:"default:false" json:"is_downloaded"`
}
