package podcasts

import (
	"fmt"
	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categoryRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/repositories"
	podcastsDto "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
	podcastsModels "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/repositories"
	"github.com/google/uuid"
	"sort"
	"strings"
)

type PodcastService struct {
	PodcastRepository *podcasts.PodcastRepository
}

func NewPodcastService(podcastRepository *podcasts.PodcastRepository) *PodcastService {
	return &PodcastService{PodcastRepository: podcastRepository}
}

func (s *PodcastService) GetAllPodcasts(getAllPodcastsRequestDto podcastsDto.GetAllPodcastsRequestDto, userID string) base.Response {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID format", err)
	}

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
		podcastDtos[i] = podcastsDto.MapToPodcastDTO(podcast, userUUID)
	}

	return base.SetPaginatedResponse(podcastDtos, page, perPage, totalCount)
}

// generateRecommendedPodcastsCacheKey creates a consistent cache key for recommended podcasts
func generateRecommendedPodcastsCacheKey(userID string, categoriesIDs []string) string {
	sortedCategories := make([]string, len(categoriesIDs))
	copy(sortedCategories, categoriesIDs)
	sort.Strings(sortedCategories)

	return fmt.Sprintf("recommended_podcasts:%s:%s", userID, strings.Join(sortedCategories, ","))
}

func (s *PodcastService) GetRecommendedPodcasts(userID string, userCategoriesIDs []string) base.Response {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID format", err)
	}

	/*	cacheKey := generateRecommendedPodcastsCacheKey(userID, userCategoriesIDs)
		ctx := context.Background()
		cachedData, err := redis.Get(ctx, cacheKey)
		if err == nil {
			var response []podcastsDto.GetRecommendedPodcastsResponseDto
			if err := json.Unmarshal([]byte(cachedData), &response); err == nil {
				return base.SetData(response)
			}
		}
	*/
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

	/*	responseJSON, err := json.Marshal(response)
		if err == nil {
			redis.SetWithTTL(ctx, cacheKey, string(responseJSON), 2*time.Hour)
		}
	*/
	return base.SetData(response)
}

func (s *PodcastService) GetTrendingPodcasts(userID string) base.Response {
	userUUID := uuid.Nil
	if userID != "" {
		parsedID, err := uuid.Parse(userID)
		if err != nil {
			return base.SetErrorMessage("Invalid user ID format", err)
		}
		userUUID = parsedID
	}

	podcasts, err := s.PodcastRepository.GetTrendingPodcasts()
	if err != nil {
		return base.SetErrorMessage("Failed to fetch trending podcasts", err)
	}

	response := make([]interface{}, len(podcasts))
	for i, podcast := range podcasts {
		response[i] = podcastsDto.MapToPodcastDTO(podcast, userUUID)
	}

	return base.SetData(response)
}

func (s *PodcastService) GetPodcastDetails(podcastID string, userID string) base.Response {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID format", err)
	}

	podcastUUID, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("Invalid podcast ID format", err)
	}

	podcast, err := s.PodcastRepository.FindPodcastByID(podcastUUID)
	if err != nil {
		return base.SetErrorMessage("Failed to get podcast details", err)
	}

	podcastDetailsDto := podcastsDto.MapToPodcastDTO(*podcast, userUUID)

	return base.SetData(podcastDetailsDto)
}

func (s *PodcastService) LikePodcast(podcastID string, addLikes int) base.Response {
	podcastUUID, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("Invalid podcast ID format", err)
	}

	if addLikes <= 0 {
		addLikes = 1
	}

	likeCount, err := s.PodcastRepository.IncrementLikesCount(podcastUUID, addLikes)
	if err != nil || likeCount == 0 {
		return base.SetErrorMessage("Failed to like podcast", err)
	}

	return base.SetData(podcastsDto.LikePodcastResponseDto{
		PodcastID:         podcastUUID.String(),
		PodcastTotalLikes: likeCount,
	})
}

func (s *PodcastService) GetPodcastsByCategory(getPodcastsByCategoryRequestDto podcastsDto.GetPodcastsByCategoryRequestDto, userID string, categoryID string) base.Response {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID", err)
	}

	categoryUUID, err := uuid.Parse(categoryID)

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
		podcastDtos[i] = podcastsDto.MapToPodcastDTO(podcast, userUUID)
	}

	return base.SetPaginatedResponse(podcastDtos, page, perPage, totalCount)
}

func (s *PodcastService) ToggleDownloadPodcast(userID, podcastID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	pid, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("Invalid Podcast ID", err)
	}

	exists, err := s.PodcastRepository.IsDownloaded(uid, pid)
	if err != nil {
		return base.SetErrorMessage("Failed to check download status", err)
	}

	var action string
	if exists {
		err = s.PodcastRepository.RemoveDownload(uid, pid)
		action = "removed"
	} else {
		err = s.PodcastRepository.AddDownload(uid, pid)
		action = "added"
	}
	if err != nil {
		return base.SetErrorMessage("Failed to toggle download", err)
	}

	return base.SetSuccessMessage("Download " + action + " successfully")
}

func (s *PodcastService) TrackUserPodcast(userID, podcastID string, trackUserPodcastRequestDto podcastsDto.TrackUserPodcastRequestDto) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID", err)
	}

	pid, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("Invalid podcast ID", err)
	}

	resumePosition := trackUserPodcastRequestDto.ResumePosition
	isCompleted := trackUserPodcastRequestDto.IsCompleted
	trackUserPodcast, err := s.PodcastRepository.TrackUserPodcast(uid, pid, resumePosition, isCompleted)
	if err != nil {
		return base.SetErrorMessage("Failed to track podcast", err)
	}

	if trackUserPodcast == nil {
		podcast, err := s.PodcastRepository.FindPodcastByID(pid)
		if err != nil {
			return base.SetErrorMessage("Failed to get podcast details", err)
		}
		if podcast == nil {
			return base.SetErrorMessage("Failed to get podcast ID", err)
		}
		userPodcast := &podcastsModels.UserPodcast{
			UserID:         uid,
			PodcastID:      podcast.ID,
			CategoryID:     podcast.CategoryID,
			ResumePosition: 0,
			IsCompleted:    false,
		}
		trackUserPodcast, err = s.PodcastRepository.CreateUserPodcast(userPodcast)
		if err != nil {
			return base.SetErrorMessage("Failed to track podcast", err)
		}
	}

	TrackUserPodcastDto := podcastsDto.TrackUserPodcastResponseDto{
		ResumePosition: trackUserPodcast.ResumePosition,
		IsCompleted:    trackUserPodcast.IsCompleted,
	}
	return base.SetData(TrackUserPodcastDto)
}

func (s *PodcastService) UserWatchHistory(userID string, getUserWatchHistoryRequestDto podcastsDto.GetUserWatchHistoryRequestDto) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID", err)
	}

	page := getUserWatchHistoryRequestDto.Page
	perPage := getUserWatchHistoryRequestDto.PerPage
	offset := (page - 1) * perPage
	limit := perPage

	userCompletedPodcasts, totalCount, err := s.PodcastRepository.GetUserCompletedPodcasts(uid, offset, limit)
	if err != nil {
		return base.SetErrorMessage("Failed to get all podcasts by IDs", err)
	}

	podcastDtos := make([]interface{}, len(*userCompletedPodcasts))
	for i, podcast := range *userCompletedPodcasts {
		podcastDtos[i] = podcastsDto.MapToPodcastDTO(podcast, uid)
	}

	return base.SetPaginatedResponse(podcastDtos, page, perPage, totalCount)
}
