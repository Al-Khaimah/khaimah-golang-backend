package podcasts

import (
	"fmt"

	podcastsModels "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
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
