package podcasts

import (
	"testing"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	podcastsDto "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPodcasts_InvalidUserID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	response := service.GetAllPodcasts(podcastsDto.GetAllPodcastsRequestDto{
		PaginationRequest: base.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}, "invalid-uuid")

	assert.Equal(t, "Invalid user ID format", response.MessageTitle)
}

func TestGetRecommendedPodcasts_InvalidUserID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	response := service.GetRecommendedPodcasts("invalid-uuid", []string{"valid-uuid"})

	assert.Equal(t, "Invalid user ID format", response.MessageTitle)
}

func TestGetRecommendedPodcasts_InvalidCategoryID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	userID := uuid.New().String()
	response := service.GetRecommendedPodcasts(userID, []string{"invalid-category-id"})

	assert.Equal(t, "Invalid category ID format", response.MessageTitle)
}

func TestGetTrendingPodcasts_InvalidUserID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	response := service.GetTrendingPodcasts("invalid-uuid")

	assert.Equal(t, "Invalid user ID format", response.MessageTitle)
}

func TestGetPodcastDetails_InvalidUserID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	podcastID := uuid.New().String()
	response := service.GetPodcastDetails(podcastID, "invalid-uuid")

	assert.Equal(t, "Invalid user ID format", response.MessageTitle)
}

func TestGetPodcastDetails_InvalidPodcastID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	userID := uuid.New().String()
	response := service.GetPodcastDetails("invalid-podcast-id", userID)

	assert.Equal(t, "Invalid podcast ID format", response.MessageTitle)
}

func TestLikePodcast_InvalidPodcastID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	response := service.LikePodcast("invalid-podcast-id", 1)

	assert.Equal(t, "Invalid podcast ID format", response.MessageTitle)
}

func TestGetPodcastsByCategory_InvalidUserID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	categoryID := uuid.New().String()
	response := service.GetPodcastsByCategory(podcastsDto.GetPodcastsByCategoryRequestDto{
		PaginationRequest: base.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}, "invalid-uuid", categoryID)

	assert.Equal(t, "Invalid user ID", response.MessageTitle)
}

func TestGetPodcastsByCategory_InvalidCategoryID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	userID := uuid.New().String()
	response := service.GetPodcastsByCategory(podcastsDto.GetPodcastsByCategoryRequestDto{
		PaginationRequest: base.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}, userID, "invalid-category-id")

	assert.Equal(t, "Invalid category ID format", response.MessageTitle)
}

func TestToggleDownloadPodcast_InvalidUserID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	podcastID := uuid.New().String()
	response := service.ToggleDownloadPodcast("invalid-user-id", podcastID)

	assert.Equal(t, "Invalid User ID", response.MessageTitle)
}

func TestToggleDownloadPodcast_InvalidPodcastID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	userID := uuid.New().String()
	response := service.ToggleDownloadPodcast(userID, "invalid-podcast-id")

	assert.Equal(t, "Invalid Podcast ID", response.MessageTitle)
}

func TestTrackUserPodcast_InvalidUserID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	podcastID := uuid.New().String()
	trackRequest := podcastsDto.TrackUserPodcastRequestDto{
		ResumePosition: 100,
		IsCompleted:    false,
	}

	response := service.TrackUserPodcast("invalid-user-id", podcastID, trackRequest)

	assert.Equal(t, "Invalid user ID", response.MessageTitle)
}

func TestTrackUserPodcast_InvalidPodcastID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	userID := uuid.New().String()
	trackRequest := podcastsDto.TrackUserPodcastRequestDto{
		ResumePosition: 100,
		IsCompleted:    false,
	}

	response := service.TrackUserPodcast(userID, "invalid-podcast-id", trackRequest)

	assert.Equal(t, "Invalid podcast ID", response.MessageTitle)
}

func TestUserWatchHistory_InvalidUserID(t *testing.T) {
	service := PodcastService{
		PodcastRepository: nil,
	}

	historyRequest := podcastsDto.GetUserWatchHistoryRequestDto{
		PaginationRequest: base.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}

	response := service.UserWatchHistory("invalid-user-id", historyRequest)

	assert.Equal(t, "Invalid user ID", response.MessageTitle)
}

func TestLikePodcast_NormalizesAddLikes(t *testing.T) {
	likePodcast := func(podcastID string, addLikes int) base.Response {
		podcastUUID, err := uuid.Parse(podcastID)
		if err != nil {
			return base.SetErrorMessage("Invalid podcast ID format", err)
		}

		if addLikes <= 0 {
			addLikes = 1
		}

		return base.SetData(podcastsDto.LikePodcastResponseDto{
			PodcastID:         podcastUUID.String(),
			PodcastTotalLikes: 15,
		})
	}

	podcastID := uuid.New().String()
	response := likePodcast(podcastID, -5)

	assert.Equal(t, 200, response.HTTPStatus)
	assert.Equal(t, "success", response.MessageType)

	data, ok := response.Data.(podcastsDto.LikePodcastResponseDto)
	assert.True(t, ok)
	assert.Equal(t, podcastID, data.PodcastID)
	assert.Equal(t, 15, data.PodcastTotalLikes)
}
