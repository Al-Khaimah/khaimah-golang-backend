package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	podcastModel "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
	userModel "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
)

type BookmarkRepository struct {
	DB *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{
		DB: db,
	}
}

func (r *BookmarkRepository) FindUserBookmarks(userID uuid.UUID) ([]podcastModel.Podcast, error) {
	var user userModel.User

	result := r.DB.Where("id = ?", userID).Preload("Bookmarks").First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return []podcastModel.Podcast{}, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return user.Bookmarks, nil
}

func (r *BookmarkRepository) IsBookmarked(userID, podcastID uuid.UUID) (bool, error) {
	var bookmark podcastModel.BookmarkPodcast
	err := r.DB.Where("user_id = ? AND podcast_id = ?", userID, podcastID).First(&bookmark).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *BookmarkRepository) AddBookmark(userID, podcastID uuid.UUID) error {
	bookmark := podcastModel.BookmarkPodcast{
		UserID:    userID,
		PodcastID: podcastID,
	}
	return r.DB.Create(&bookmark).Error
}

func (r *BookmarkRepository) RemoveBookmark(userID, podcastID uuid.UUID) error {
	return r.DB.Where("user_id = ? AND podcast_id = ?", userID, podcastID).Delete(&podcastModel.BookmarkPodcast{}).Error
}
