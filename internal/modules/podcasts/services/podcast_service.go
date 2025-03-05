package podcasts

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	podcastsDto "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PodcastService struct {
	PodcastRepository *podcasts.PodcastRepository
}

func NewPodcastService(podcastRepository *podcasts.PodcastRepository) *PodcastService {
	return &PodcastService{PodcastRepository: podcastRepository}
}

func (s *PodcastService) GetAllPodcasts(
	c echo.Context,
	getAllPodcastsRequestDto podcastsDto.GetAllPodcastsRequestDto,
) (base.Response, error) {
	page := getAllPodcastsRequestDto.Page
	perPage := getAllPodcastsRequestDto.PerPage

	offset := (page - 1) * perPage
	limit := perPage

	podcasts, totalCount, err := s.PodcastRepository.GetAllPodcasts(offset, limit)
	if err != nil {
		return base.SetErrorMessage("Failed to get podcasts", err), err
	}
	podcastDtos := make([]interface{}, len(podcasts))
	for i, podcast := range podcasts {
		podcastDtos[i] = podcastsDto.PodcastDto{
			ID:                    podcast.ID.String(),
			Title:                 podcast.Title,
			Description:           podcast.Description,
			AudioURL:              podcast.AudioURL,
			CoverImageURL:         podcast.CoverImageURL,
			CoverImageDescription: podcast.CoverImageDescription,
			LikesCount:            podcast.LikesCount,
			CategoryID:            podcast.CategoryID.String(),
		}
	}
	return base.SetPaginatedResponse(podcastDtos, page, perPage, totalCount), nil
}

func (s *PodcastService) GetRecommendedPodcasts(userID string, userCategoriesIds []string) (base.Response, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID format", err), err
	}

	categoriesUUID := make([]uuid.UUID, len(userCategoriesIds))
	for i, categoryId := range userCategoriesIds {
		categoriesUUID[i], err = uuid.Parse(categoryId)
		if err != nil {
			return base.SetErrorMessage("Invalid category ID format", err), err
		}
	}

	listenedPodcastIDs, err := s.PodcastRepository.GetlistenedPodcastIDs(userUUID)
	if err != nil {
		return base.SetErrorMessage("Failed to get listened podcast IDs", err), err
	}

	recommendedPodcasts, err := s.PodcastRepository.GetRecommendedPodcasts(listenedPodcastIDs, categoriesUUID)
	if err != nil {
		return base.SetErrorMessage("Failed to get recommended podcasts", err), err
	}

	recommendedPodcastsResponseDto := make([]podcastsDto.PodcastDto, len(recommendedPodcasts))
	for i, podcast := range recommendedPodcasts {
		recommendedPodcastsResponseDto[i] = podcastsDto.PodcastDto{
			ID:                    podcast.ID.String(),
			Title:                 podcast.Title,
			Description:           podcast.Description,
			AudioURL:              podcast.AudioURL,
			CoverImageURL:         podcast.CoverImageURL,
			CoverImageDescription: podcast.CoverImageDescription,
			LikesCount:            podcast.LikesCount,
			CategoryID:            podcast.CategoryID.String(),
		}
	}
	return base.SetData(recommendedPodcastsResponseDto), nil
}
