package podcasts

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	podcastsModels "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
	podcastRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/repositories"
	"github.com/google/uuid"
)

type PodcastDto struct {
	ID                    string `json:"id"`
	Title                 string `json:"title"`
	Content               string `json:"content,omitempty"`
	AudioURL              string `json:"audio_url"`
	CoverImageURL         string `json:"cover_image_url"`
	CoverImageDescription string `json:"cover_image_description,omitempty"`
	LikesCount            int    `json:"likes_count"`
	Duration              int    `json:"duration"`
	CategoryID            string `json:"category_id"`
	IsDownloaded          bool   `json:"is_downloaded"`
	IsBookmarked          bool   `json:"is_bookmarked"`
	IsCompleted           bool   `json:"is_completed"`
	IsTrending            bool   `json:"is_trending"`
	CreatedAt             string `json:"created_at,omitempty"`
	UpdatedAt             string `json:"updated_at,omitempty"`
	DeletedAt             string `json:"deleted_at,omitempty"`
}

func MapToPodcastDTO(podcast podcastsModels.Podcast, userID uuid.UUID) PodcastDto {
	var isDownloaded, isBookmarked, isCompleted, isTrending bool

	r := podcastRepository.NewPodcastRepository(config.GetDB())
	if userID != uuid.Nil {
		isDownloaded, _ = r.IsDownloaded(userID, podcast.ID)
		isBookmarked, _ = r.IsBookmarked(userID, podcast.ID)
		isCompleted, _ = r.IsCompleted(userID, podcast.ID)
	}
	isTrending, _ = r.IsTrending(podcast.ID)

	return PodcastDto{
		ID:            podcast.ID.String(),
		Title:         podcast.Title,
		AudioURL:      podcast.AudioURL,
		CoverImageURL: podcast.CoverImageURL,
		LikesCount:    podcast.LikesCount,
		Duration:      podcast.Duration,
		CategoryID:    podcast.CategoryID.String(),
		IsDownloaded:  isDownloaded,
		IsBookmarked:  isBookmarked,
		IsCompleted:   isCompleted,
		IsTrending:    isTrending,
	}
}

type GetAllPodcastsRequestDto struct {
	base.PaginationRequest
}

type GetRecommendedPodcastsResponseDto struct {
	CategoryID   string       `json:"category_id"`
	CategoryName string       `json:"category_name"`
	Podcasts     []PodcastDto `json:"podcasts"`
}

type GetPodcastDetailsRequestDto struct {
	ID string `json:"id" param:"id" validate:"required,uuid" message:"ID must be a valid ID format"`
}

type LikePodcastRequestDto struct {
	AddLikes int `json:"add_likes" validate:"omitempty,min=1"`
}

type LikePodcastResponseDto struct {
	PodcastID         string `json:"podcast_id"`
	PodcastTotalLikes int    `json:"podcast_total_likes"`
}

type GetPodcastsByCategoryRequestDto struct {
	base.PaginationRequest
}

type GetPodcastsByCategoryResponseDto struct {
	Podcasts []PodcastDto `json:"podcasts"`
}

type TrackUserPodcastRequestDto struct {
	ResumePosition int  `json:"resume_position"`
	IsCompleted    bool `json:"is_completed"`
}

type TrackUserPodcastResponseDto struct {
	ResumePosition int  `json:"resume_position"`
	IsCompleted    bool `json:"is_completed"`
}

type GetUserWatchHistoryRequestDto struct {
	base.PaginationRequest
}
