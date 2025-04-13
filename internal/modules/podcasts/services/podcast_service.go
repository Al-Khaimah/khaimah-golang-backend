package podcasts

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categoryRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/repositories"
	podcastsDto "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/repositories"
	"github.com/google/uuid"
)

type PodcastService struct {
	PodcastRepository *podcasts.PodcastRepository
}

func NewPodcastService(podcastRepository *podcasts.PodcastRepository) *PodcastService {
	return &PodcastService{PodcastRepository: podcastRepository}
}

func (s *PodcastService) GetAllPodcasts(getAllPodcastsRequestDto podcastsDto.GetAllPodcastsRequestDto) base.Response {
	page := getAllPodcastsRequestDto.Page
	perPage := getAllPodcastsRequestDto.PerPage

	offset := (page - 1) * perPage
	limit := perPage

	podcasts, totalCount, err := s.PodcastRepository.GetAllPodcasts(offset, limit)
	if err != nil {
		return base.SetErrorMessage("Failed to get podcasts", err)
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
	return base.SetPaginatedResponse(podcastDtos, page, perPage, totalCount)
}

func (s *PodcastService) GetRecommendedPodcasts(userID string, userCategoriesIDs []string) base.Response {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID format", err)
	}

	var categoriesUUID []uuid.UUID
	for _, id := range userCategoriesIDs {
		catID, err := uuid.Parse(id)
		if err != nil {
			return base.SetErrorMessage("Invalid category ID format", err)
		}
		categoriesUUID = append(categoriesUUID, catID)
	}

	completedIDs, err := s.PodcastRepository.GetCompletedPodcastsIDs(userUUID)
	if err != nil {
		return base.SetErrorMessage("Failed to get completed podcast IDs", err)
	}

	podcasts, err := s.PodcastRepository.GetRecommendedPodcasts(categoriesUUID, completedIDs)
	if err != nil {
		return base.SetErrorMessage("Failed to get recommended podcasts", err)
	}

	categoryRepo := categoryRepository.NewCategoryRepository(config.GetDB())
	grouped := make(map[string]podcastsDto.GetRecommendedPodcastsResponseDto)

	for _, podcast := range podcasts {
		dto := podcastsDto.MapToPodcastDTO(podcast, userUUID)
		categoryID := podcast.CategoryID.String()

		if _, exists := grouped[categoryID]; !exists {
			category, _ := categoryRepo.FindCategoryByID(podcast.CategoryID)
			grouped[categoryID] = podcastsDto.GetRecommendedPodcastsResponseDto{
				CategoryID:   categoryID,
				CategoryName: category.Name,
				Podcasts:     []podcastsDto.PodcastDto{},
			}
		}

		group := grouped[categoryID]
		group.Podcasts = append(group.Podcasts, dto)
		grouped[categoryID] = group
	}

	var response []podcastsDto.GetRecommendedPodcastsResponseDto
	for _, group := range grouped {
		response = append(response, group)
	}

	return base.SetData(response)
}

func (s *PodcastService) GetPodcastDetails(podcastID string) base.Response {
	podcastUUID, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("Invalid podcast ID format", err)
	}

	podcast, err := s.PodcastRepository.FindPodcastByID(podcastUUID)
	if err != nil {
		return base.SetErrorMessage("Failed to get podcast details", err)
	}

	podcastDetailsDto := podcastsDto.PodcastDto{
		ID:                    podcast.ID.String(),
		Title:                 podcast.Title,
		Description:           podcast.Description,
		AudioURL:              podcast.AudioURL,
		CoverImageURL:         podcast.CoverImageURL,
		CoverImageDescription: podcast.CoverImageDescription,
		LikesCount:            podcast.LikesCount,
		CategoryID:            podcast.CategoryID.String(),
	}

	return base.SetData(podcastDetailsDto)
}

func (s *PodcastService) LikePodcast(podcastID string) base.Response {
	podcastUUID, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("Invalid podcast ID format", err)
	}

	likeCount, err := s.PodcastRepository.IncrementLikesCount(podcastUUID)
	if err != nil || likeCount == 0 {
		return base.SetErrorMessage("Failed to like podcast", err)
	}

	likePodcastResponseDto := podcastsDto.LikePodcastResponseDto{
		PodcastID:         podcastUUID.String(),
		PodcastTotalLikes: likeCount,
	}

	return base.SetData(likePodcastResponseDto)
}

func (s *PodcastService) GetPodcastsByCategory(getPodcastsByCategoryRequestDto podcastsDto.GetPodcastsByCategoryRequestDto) base.Response {
	categoryUUID, err := uuid.Parse(getPodcastsByCategoryRequestDto.CategoryID)

	page := getPodcastsByCategoryRequestDto.Page
	perPage := getPodcastsByCategoryRequestDto.PerPage

	offset := (page - 1) * perPage
	limit := perPage

	if err != nil {
		return base.SetErrorMessage("Invalid category ID format", err)
	}

	podcasts, totalCount, err := s.PodcastRepository.FindPodcastsByCategoryID(categoryUUID, offset, limit)
	if err != nil {
		return base.SetErrorMessage("Failed to get podcasts by category ID", err)
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

	return base.SetPaginatedResponse(podcastDtos, page, perPage, totalCount)
}

func (s *PodcastService) DownloadPodcast(userID, podcastID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID", err)
	}

	pid, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("Invalid podcast ID", err)
	}

	err = s.PodcastRepository.MarkPodcastAsDownloaded(uid, pid)
	if err != nil {
		return base.SetErrorMessage("Failed to mark podcast as downloaded", err)
	}

	return base.SetSuccessMessage("Podcast marked as downloaded")
}
