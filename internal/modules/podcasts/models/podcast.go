package podcasts

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/google/uuid"
)

type FetchedFrom string

const (
	FetchedFromGrokAPI      FetchedFrom = "grokApi"
	FetchedFromWorldNewsAPI FetchedFrom = "worldNewsApi"
)

type Podcast struct {
	base.Model
	Title                 string      `gorm:"type:varchar(255);index" json:"title"`
	Content               string      `gorm:"type:text" json:"content"`
	AudioURL              string      `gorm:"type:text" json:"audio_url"`
	CoverImageURL         string      `gorm:"type:text" json:"cover_image_url"`
	CoverImageDescription string      `gorm:"type:text" json:"cover_image_description"`
	LikesCount            int         `gorm:"default:0" json:"likes_count"`
	Duration              int         `gorm:"default:0" json:"duration"`
	CategoryID            uuid.UUID   `gorm:"type:uuid;index" json:"category_id"`
	FetchedFrom           FetchedFrom `gorm:"type:text" json:"fetched_from"`
	Tags                  string      `gorm:"type:text" json:"tags"`
}

type UserPodcast struct {
	base.Model
	UserID         uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	PodcastID      uuid.UUID `gorm:"type:uuid;index" json:"podcast_id"`
	CategoryID     uuid.UUID `gorm:"type:uuid;index" json:"category_id"`
	ResumePosition int       `gorm:"default:0" json:"resume_position"`
	IsCompleted    bool      `gorm:"default:false" json:"is_completed"`
}

type BookmarkPodcast struct {
	UserID    uuid.UUID `gorm:"primaryKey;type:uuid"`
	PodcastID uuid.UUID `gorm:"primaryKey;type:uuid"`
}

func (BookmarkPodcast) TableName() string {
	return "user_bookmarks" //by default it will be named 'user_bookmarks'
}
