package users

import (
	"testing"

	userDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestUpdateUserProfile_InvalidUserID tests validation of user ID
func TestUpdateUserProfile_InvalidUserID(t *testing.T) {
	service := UserService{}
	response := service.UpdateUserProfile("invalid-user-id", userDTO.UpdateProfileDTO{
		FirstName: "John",
		LastName:  "Doe",
	})
	assert.Equal(t, "Invalid User ID", response.MessageTitle)
}

// TestChangePassword_InvalidUserID tests validation of user ID
func TestChangePassword_InvalidUserID(t *testing.T) {
	service := UserService{}
	response := service.ChangePassword("invalid-user-id", userDTO.ChangePasswordDTO{
		OldPassword: "oldPassword",
		NewPassword: "newPassword",
	})
	assert.Equal(t, "Invalid User ID", response.MessageTitle)
}

// TestDeleteUser_InvalidUserID tests validation of user ID
func TestDeleteUser_InvalidUserID(t *testing.T) {
	service := UserService{}
	response := service.DeleteUser("invalid-user-id")
	assert.Equal(t, "Invalid User ID", response.MessageTitle)
}

// TestGetUserProfile_InvalidUserID tests validation of user ID
func TestGetUserProfile_InvalidUserID(t *testing.T) {
	service := UserService{}
	response := service.GetUserProfile("invalid-user-id")
	assert.Equal(t, "Invalid User ID", response.MessageTitle)
}

// TestUpdateUserPreferences_InvalidUserID tests validation of user ID
func TestUpdateUserPreferences_InvalidUserID(t *testing.T) {
	service := UserService{}
	_, err := service.UpdateUserPreferences("invalid-user-id", userDTO.UpdatePreferencesDTO{
		Categories: []string{"category1", "category2"},
	})
	assert.Equal(t, "الرقم التعريفي للمستخدم غير صالح", err.Error())
}

// TestGetUserBookmarks_InvalidUserID tests validation of user ID
func TestGetUserBookmarks_InvalidUserID(t *testing.T) {
	service := UserService{}
	response := service.GetUserBookmarks("invalid-user-id")
	assert.Equal(t, "Invalid User ID", response.MessageTitle)
}

// TestToggleBookmarkPodcast_InvalidUserID tests validation of user ID
func TestToggleBookmarkPodcast_InvalidUserID(t *testing.T) {
	service := UserService{}
	podcastID := uuid.New().String()
	response := service.ToggleBookmarkPodcast("invalid-user-id", podcastID)
	assert.Equal(t, "Invalid User ID", response.MessageTitle)
}

// TestToggleBookmarkPodcast_InvalidPodcastID tests validation of podcast ID
func TestToggleBookmarkPodcast_InvalidPodcastID(t *testing.T) {
	service := UserService{}
	userID := uuid.New().String()
	response := service.ToggleBookmarkPodcast(userID, "invalid-podcast-id")
	assert.Equal(t, "Invalid Podcast ID", response.MessageTitle)
}

// TestGetDownloadedPodcasts_InvalidUserID tests validation of user ID
func TestGetDownloadedPodcasts_InvalidUserID(t *testing.T) {
	service := UserService{}
	response := service.GetDownloadedPodcasts("invalid-user-id")
	assert.Equal(t, "Invalid User ID", response.MessageTitle)
}
