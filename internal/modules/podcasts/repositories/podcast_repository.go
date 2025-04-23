package podcasts

import (
	"fmt"
	"time"

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

func (r *PodcastRepository) GetCompletedPodcastsIDs(userUUID uuid.UUID) ([]uuid.UUID, error) {
	var completedPodcastsIDs []uuid.UUID
	result := r.DB.Model(&podcastsModels.UserPodcast{}).
		Where("user_id = ?", userUUID).
		Where("is_completed = ?", true).
		Pluck("podcast_id", &completedPodcastsIDs)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get listened podcast IDs: %w", result.Error)
	}
	return completedPodcastsIDs, nil
}

func (r *PodcastRepository) GetRecommendedPodcasts(categoriesUUID []uuid.UUID, completedPodcastsIDs []uuid.UUID) ([]podcastsModels.Podcast, error) {
	var podcasts []podcastsModels.Podcast

	subQuery := r.DB.Model(&podcastsModels.Podcast{}).
		Select("*, ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY created_at DESC) as rn").
		Where("category_id IN ?", categoriesUUID)

	if len(completedPodcastsIDs) > 0 {
		subQuery = subQuery.Where("id NOT IN ?", completedPodcastsIDs)
	}

	result := r.DB.Table("(?) as p", subQuery).
		Where("p.rn <= ?", 6).
		Find(&podcasts)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recommended podcasts: %w", result.Error)
	}

	return podcasts, nil
}

func (r *PodcastRepository) GetTrendingPodcasts() ([]podcastsModels.Podcast, error) {
	var podcasts []podcastsModels.Podcast
	tenDaysAgo := time.Now().AddDate(0, 0, -10)

	result := r.DB.Model(&podcastsModels.Podcast{}).
		Where("created_at >= ?", tenDaysAgo).
		Order("likes_count DESC").
		Limit(10).
		Find(&podcasts)

	if result.Error != nil {
		return nil, result.Error
	}

	return podcasts, nil
}

func (r *PodcastRepository) FindPodcastByID(podcastID uuid.UUID) (podcastsModels.Podcast, error) {
	var podcast podcastsModels.Podcast
	result := r.DB.Where("id = ?", podcastID).First(&podcast)

	if result.Error == gorm.ErrRecordNotFound {
		return podcastsModels.Podcast{}, nil
	}
	if result.Error != nil {
		return podcastsModels.Podcast{}, fmt.Errorf("failed to find podcast: %w", result.Error)
	}

	return podcast, nil
}

func (r *PodcastRepository) IncrementLikesCount(podcastID uuid.UUID, addLikes int) (int, error) {
	var podcast podcastsModels.Podcast

	result := r.DB.Model(&podcastsModels.Podcast{}).
		Where("id = ?", podcastID).
		UpdateColumn("likes_count", gorm.Expr("likes_count + ?", addLikes)).
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
