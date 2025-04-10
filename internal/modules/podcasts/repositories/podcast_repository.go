package podcasts

import (
	"fmt"
	podcastDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
	podcastsModels "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PodcastRepository struct {
	DB *gorm.DB
}

func NewPodcastRepository(DB *gorm.DB) *PodcastRepository {
	return &PodcastRepository{DB: DB}
}

func (r *PodcastRepository) GetAllPodcasts(offset int, limit int) ([]podcastsModels.Podcast, int, error) {
	var podcasts []podcastsModels.Podcast
	var totalCount int64

	result := r.DB.Model(&podcastsModels.Podcast{}).Count(&totalCount)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to get all podcasts: %w", result.Error)
	}
	result = r.DB.Limit(limit).Offset(offset).Find(&podcasts)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to get all podcasts: %w", result.Error)
	}

	return podcasts, int(totalCount), nil
}

func (r *PodcastRepository) MapToPodcastDTO(podcast podcastsModels.Podcast, userID uuid.UUID) podcastDTO.PodcastDto {
	isDownloaded, _ := r.IsDownloaded(userID, podcast.ID)
	isBookmarked, _ := r.IsBookmarked(userID, podcast.ID)
	isCompleted, _ := r.IsCompleted(userID, podcast.ID)

	return podcastDTO.PodcastDto{
		ID:                    podcast.ID.String(),
		Title:                 podcast.Title,
		Description:           podcast.Description,
		AudioURL:              podcast.AudioURL,
		CoverImageURL:         podcast.CoverImageURL,
		CoverImageDescription: podcast.CoverImageDescription,
		LikesCount:            podcast.LikesCount,
		CategoryID:            podcast.CategoryID.String(),
		IsDownloaded:          isDownloaded,
		IsBookmarked:          isBookmarked,
		IsCompleted:           isCompleted,
	}
}

func (r *PodcastRepository) GetListenedPodcastIDs(userUUID uuid.UUID) ([]uuid.UUID, error) {
	var listenedPodcastIDs []uuid.UUID
	result := r.DB.Model(&podcastsModels.UserPodcast{}).
		Where("user_id = ?", userUUID).
		Pluck("podcast_id", &listenedPodcastIDs)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get listened podcast IDs: %w", result.Error)
	}
	return listenedPodcastIDs, nil
}

func (r *PodcastRepository) GetRecommendedPodcasts(listenedPodcastIDs []uuid.UUID, categoriesUUID []uuid.UUID) ([]podcastsModels.Podcast, error) {
	var podcasts []podcastsModels.Podcast

	result := r.DB.Model(&podcastsModels.Podcast{}).
		Where("category_id IN ?", categoriesUUID).
		Where("id NOT IN ?", listenedPodcastIDs).
		Order("created_at DESC").
		Limit(10).
		Find(&podcasts)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recommended podcasts: %w", result.Error)
	}
	return podcasts, nil
}

func (r *PodcastRepository) FindPodcastByID(podcastID uuid.UUID) (*podcastsModels.Podcast, error) {
	var podcast podcastsModels.Podcast
	result := r.DB.Where("id = ?", podcastID).First(&podcast)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find podcast: %w", result.Error)
	}

	return &podcast, nil
}

func (r *PodcastRepository) IncrementLikesCount(podcastID uuid.UUID) (int, error) {
	var podcast podcastsModels.Podcast
	result := r.DB.Model(&podcastsModels.Podcast{}).
		Where("id = ?", podcastID).
		Update("likes_count", gorm.Expr("likes_count + 1")).
		First(&podcast)

	if result.Error == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if result.Error != nil {
		return 0, fmt.Errorf("failed to increment likes count: %w", result.Error)
	}

	return podcast.LikesCount, nil
}

func (r *PodcastRepository) FindPodcastsByCategoryID(categoryID uuid.UUID, offset int, limit int) ([]podcastsModels.Podcast, int, error) {
	var podcasts []podcastsModels.Podcast
	var totalCount int64

	result := r.DB.Model(&podcastsModels.Podcast{}).
		Where("category_id = ?", categoryID).
		Offset(offset).
		Limit(limit).
		Find(&podcasts)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find podcasts by category ID: %w", result.Error)
	}

	result = r.DB.Model(&podcastsModels.Podcast{}).
		Where("category_id = ?", categoryID).
		Count(&totalCount)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to find podcasts by category ID: %w", result.Error)
	}

	return podcasts, int(totalCount), nil
}

func (r *PodcastRepository) MarkPodcastAsDownloaded(userID, podcastID uuid.UUID) error {
	var user users.User
	if err := r.DB.Preload("Downloads").First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	podcast, err := r.FindPodcastByID(podcastID)
	if err != nil {
		return err
	}

	var existingPodcasts []podcastsModels.Podcast
	err = r.DB.Model(&user).Association("Downloads").Find(&existingPodcasts, "id = ?", podcastID)
	if err == nil && len(existingPodcasts) > 0 {
		return nil
	}

	return r.DB.Model(&user).Association("Downloads").Append(podcast)
}

func (r *PodcastRepository) IsDownloaded(userID, podcastID uuid.UUID) (bool, error) {
	var user users.User
	err := r.DB.Preload("Downloads", "podcasts.id = ?", podcastID).First(&user, "id = ?", userID).Error
	if err != nil {
		return false, err
	}
	return len(user.Downloads) > 0, nil
}

func (r *PodcastRepository) IsBookmarked(userID, podcastID uuid.UUID) (bool, error) {
	var user users.User
	err := r.DB.Preload("Bookmarks", "podcasts.id = ?", podcastID).First(&user, "id = ?", userID).Error
	if err != nil {
		return false, err
	}
	return len(user.Bookmarks) > 0, nil
}

func (r *PodcastRepository) IsCompleted(userID, podcastID uuid.UUID) (bool, error) {
	var userPodcast podcastsModels.UserPodcast
	err := r.DB.Where("user_id = ? AND podcast_id = ? AND is_completed = ?", userID, podcastID, true).First(&userPodcast).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
