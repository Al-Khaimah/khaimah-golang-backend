package podcasts

import (
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
		return nil, 0, result.Error
	}
	result = r.DB.Limit(limit).Offset(offset).Find(&podcasts)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return podcasts, int(totalCount), nil
}

func (r *PodcastRepository) GetlistenedPodcastIDs(userUUID uuid.UUID) ([]uuid.UUID, error) {
	var listenedPodcastIDs []uuid.UUID
	result := r.DB.Model(&podcastsModels.UserPodcast{}).
		Where("user_id = ?", userUUID).
		Pluck("podcast_id", &listenedPodcastIDs)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return listenedPodcastIDs, nil
}

func (r *PodcastRepository) GetRecommendedPodcasts(
	listenedPodcastIDs []uuid.UUID,
	categoriesUUID []uuid.UUID,
) ([]podcastsModels.Podcast, error) {
	var podcasts []podcastsModels.Podcast

	result := r.DB.Model(&podcastsModels.Podcast{}).
		Where("category_id IN ?", categoriesUUID).
		Where("id NOT IN ?", listenedPodcastIDs).
		Order("created_at DESC").
		Limit(10).
		Find(&podcasts)

	if result.Error != nil {
		return nil, result.Error
	}
	return podcasts, nil
}
