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
	CreatedAt             string `json:"created_at,omitempty"`
	UpdatedAt             string `json:"updated_at,omitempty"`
}

type GetAllPodcastsRequestDto struct {
	base.PaginationRequest
}

type GetRecommendedPodcastsResponseDto struct {
	Podcasts []PodcastDto `json:"podcasts"`
}
