package podcasts

import "github.com/Al-Khaimah/khaimah-golang-backend/internal/base"

type PodcastDto struct {
	ID                    string `json:"id"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
	AudioURL              string `json:"audio_url"`
	CoverImageURL         string `json:"cover_image_url"`
	CoverImageDescription string `json:"cover_image_description"`
	LikesCount            int    `json:"likes_count"`
	CategoryID            string `json:"category_id"`
	IsDownloaded          bool   `json:"is_downloaded"`
	IsBookmarked          bool   `json:"is_bookmarked"`
	IsCompleted           bool   `json:"is_completed"`
	CreatedAt             string `json:"created_at,omitempty"`
	UpdatedAt             string `json:"updated_at,omitempty"`
	DeletedAt             string `json:"deleted_at,omitempty"`
}

type GetAllPodcastsRequestDto struct {
	base.PaginationRequest
}

type GetRecommendedPodcastsResponseDto struct {
	Podcasts []PodcastDto `json:"podcasts"`
}

type GetPodcastDetailsRequestDto struct {
	ID string `json:"id" param:"id" validate:"required,uuid" message:"ID must be a valid ID format"`
}

type LikePodcastRequestDto struct {
	ID string `json:"id" param:"id" validate:"required,uuid" message:"ID must be a valid ID format"`
}

type LikePodcastResponseDto struct {
	PodcastID         string `json:"podcast_id"`
	PodcastTotalLikes int    `json:"podcast_total_likes"`
}

type GetPodcastsByCategoryRequestDto struct {
	base.PaginationRequest
	CategoryID string `json:"category_id" param:"category_id" validate:"required,uuid" message:"Category ID must be a valid ID format"`
}

type GetPodcastsByCategoryResponseDto struct {
	Podcasts []PodcastDto `json:"podcasts"`
}
