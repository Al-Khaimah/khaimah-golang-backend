package users

import (
	"fmt"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base/utils"
	userDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/enums"
	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	repos "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	podcastDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
)

type UserService struct {
	UserRepo      *repos.UserRepository
	AuthRepo      *repos.AuthRepository
	BookmarksRepo *repos.BookmarkRepository
}

func NewUserService(userRepo *repos.UserRepository, authRepo *repos.AuthRepository, bookmarksRepo *repos.BookmarkRepository) *UserService {
	return &UserService{
		UserRepo:      userRepo,
		AuthRepo:      authRepo,
		BookmarksRepo: bookmarksRepo,
	}
}

func (s *UserService) CreateUser(user *userDTO.SignupRequestDTO) base.Response {
	existingUser, err := s.UserRepo.FindOneByEmail(user.Email)
	if err != nil {
		return base.SetErrorMessage("Database error", err)
	}
	if existingUser != nil {
		return base.SetErrorMessage("This email is already in use", "User already exists")
	}

	categories := utils.ConvertIDsToCategories(user.Categories)
	newUser := &models.User{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      utils.FormatEmail(user.Email),
		Categories: categories,
	}

	createdUser, err := s.UserRepo.CreateUser(newUser)
	if err != nil {
		return base.SetErrorMessage("Failed to create user", err)
	}

	if createdUser == nil {
		return base.SetErrorMessage("Failed to create user", "User creation returned nil")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return base.SetErrorMessage("Failed to hash password", err)
	}

	newUserAuth := &models.IamAuth{
		UserID:   createdUser.ID,
		Password: string(hashedPassword),
		IsActive: true,
	}

	if err := s.AuthRepo.CreateUserAuth(newUserAuth); err != nil {
		return base.SetErrorMessage("Failed to create user authentication", err)
	}

	token, err := generateJWT(createdUser)
	if err != nil {
		return base.SetErrorMessage("Failed to generate token", err)
	}
	userResponse := userDTO.SignupResponseDTO{
		ID:         createdUser.ID.String(),
		FirstName:  createdUser.FirstName,
		LastName:   createdUser.LastName,
		Email:      createdUser.Email,
		Categories: createdUser.Categories,
		Token:      token,
		ExpiresAt:  "never",
	}

	slackMessage := fmt.Sprintf("🚀 New user account created:\n%s (%s)", createdUser.FirstName, createdUser.Email)
	_ = utils.SendSlackNotification(slackMessage)

	return base.SetData(userResponse, "Account created successfully")
}

func (s *UserService) LoginUser(user *userDTO.LoginRequestDTO) base.Response {
	existingUser, err := s.UserRepo.FindOneByEmail(user.Email)
	if err != nil {
		return base.SetErrorMessage("Database error", err)
	}
	if existingUser == nil {
		return base.SetErrorMessage("Invalid credentials", "Email not found")
	}

	userAuth, err := s.AuthRepo.FindAuthByUserID(existingUser.ID)
	if err != nil {
		return base.SetErrorMessage("Authentication error", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(user.Password)); err != nil {
		return base.SetErrorMessage("Invalid credentials", "Incorrect password")
	}

	token, err := generateJWT(existingUser)
	if err != nil {
		return base.SetErrorMessage("Failed to generate token", err)
	}

	userAuth.IsActive = true
	if err := s.AuthRepo.UpdateAuth(userAuth); err != nil {
		return base.SetErrorMessage("Failed to update authentication", err)
	}

	loginResponse := userDTO.LoginResponseDTO{
		ID:         existingUser.ID.String(),
		FirstName:  existingUser.FirstName,
		LastName:   existingUser.LastName,
		Email:      existingUser.Email,
		Categories: existingUser.Categories,
		Token:      token,
		ExpiresAt:  "never",
	}

	return base.SetData(loginResponse, "Logged in successfully")
}

func generateJWT(user *models.User) (string, error) {
	jwtSecret := config.GetEnv("JWT_SECRET", "alkhaimah123")
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   utils.FormatEmail(user.Email),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (s *UserService) LogoutUser(c echo.Context) base.Response {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return base.SetErrorMessage("Unauthorized", "No token provided")
	}

	userID, err := utils.ExtractUserIDFromToken(token)
	if err != nil {
		return base.SetErrorMessage("Invalid token", err)
	}

	authRecord, err := s.AuthRepo.FindAuthByUserID(userID)
	if err != nil {
		return base.SetErrorMessage("Failed to find user authentication", err)
	}

	authRecord.IsActive = false
	err = s.AuthRepo.UpdateAuth(authRecord)
	if err != nil {
		return base.SetErrorMessage("Failed to logout", err)
	}

	return base.SetSuccessMessage("Successfully logged out")
}

func (s *UserService) GetUserProfile(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("Failed to fetch user profile", err)
	}
	if user == nil {
		return base.SetErrorMessage("User not found", "No user exists with this ID")
	}

	profileResponse := userDTO.UserProfileDTO{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Categories: user.Categories,
	}

	return base.SetData(profileResponse)
}

func (s *UserService) UpdateUserProfile(userID string, updateData userDTO.UpdateProfileDTO) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("Database error", err)
	}
	if user == nil {
		return base.SetErrorMessage("User not found", "No user exists with this ID")
	}

	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}

	err = s.UserRepo.UpdateUser(user)
	if err != nil {
		return base.SetErrorMessage("Failed to update profile", err)
	}

	profileResponse := userDTO.UserProfileDTO{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return base.SetData(profileResponse, "Profile updated successfully")
}

func (s *UserService) UpdateUserPreferences(userID string, updateData userDTO.UpdatePreferencesDTO) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("Database error", err)
	}
	if user == nil {
		return base.SetErrorMessage("User not found", "No user exists with this ID")
	}

	newCategories := utils.ConvertIDsToCategories(updateData.Categories)

	err = s.UserRepo.UpdateUserPreferences(user, newCategories)
	if err != nil {
		return base.SetErrorMessage("Failed to update preferences", err)
	}

	preferencesResponse := userDTO.UpdatePreferencesDTO{
		Categories: utils.ConvertCategoriesToString(user.Categories),
	}

	return base.SetData(preferencesResponse, "User preferences updated successfully")
}

func (s *UserService) ChangePassword(userID string, req userDTO.ChangePasswordDTO) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	userAuth, err := s.AuthRepo.FindAuthByUserID(uid)
	if err != nil {
		return base.SetErrorMessage("Database error", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(req.OldPassword))
	if err != nil {
		return base.SetErrorMessage("Invalid credentials", "Old password is incorrect")
	}

	if req.OldPassword == req.NewPassword {
		return base.SetErrorMessage("Invalid request", "New password must be different from old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return base.SetErrorMessage("Failed to hash password", err)
	}

	userAuth.Password = string(hashedPassword)

	err = s.AuthRepo.UpdateAuth(userAuth)
	if err != nil {
		return base.SetErrorMessage("Failed to change password", err)
	}

	return base.SetSuccessMessage("Password changed successfully")
}

func (s *UserService) GetAllUsers(c echo.Context) base.Response {
	users, err := s.UserRepo.FindAllUsers()
	if err != nil {
		return base.SetErrorMessage("Failed to fetch users", err)
	}

	var userResponses []interface{}
	for _, user := range users {
		userResponses = append(userResponses, userDTO.UserProfileDTO{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}
	return base.SetDataPaginated(c, userResponses)
}

func (s *UserService) MarkUserAdmin(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid user ID", err)
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("Failed to fetch user", err)
	}
	if user == nil {
		return base.SetErrorMessage("User not found", nil)
	}

	user.UserType = users.UserTypeAdmin
	if err := s.UserRepo.UpdateUser(user); err != nil {
		return base.SetErrorMessage("Failed to update user", err)
	}

	return base.SetSuccessMessage("User marked as admin successfully")
}

func (s *UserService) DeleteUser(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("Database error", err)
	}
	if user == nil {
		return base.SetErrorMessage("User not found", "No user exists with this ID")
	}

	if user.UserType == users.UserTypeAdmin {
		return base.SetErrorMessage("Admin User", "You cant delete an Admin User")
	}

	if user.DeletedAt.Valid {
		return base.SetErrorMessage("User already deleted", "This user account has already been removed")
	}

	err = s.UserRepo.DeleteUser(uid)
	if err != nil {
		return base.SetErrorMessage("Failed to delete user", err)
	}

	return base.SetSuccessMessage("User deleted successfully")
}

func (s *UserService) GetUserCategoriesIDs(userID string) ([]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	userCategories, err := s.UserRepo.FindUserCategories(uid)
	if err != nil {
		return nil, err
	}
	if userCategories == nil {
		return nil, nil
	}

	categoryIDs := make([]string, len(userCategories))
	for i, userCategories := range userCategories {
		categoryIDs[i] = userCategories.ID.String()
	}

	return categoryIDs, nil
}

func (s *UserService) GetUserBookmarks(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	bookmarks, err := s.BookmarksRepo.FindUserBookmarks(uid)
	if err != nil {
		return base.SetErrorMessage("Failed to fetch bookmarks", err)
	}

	bookmarksResponse := make([]interface{}, len(bookmarks))
	for i, podcast := range bookmarks {
		bookmarksResponse[i] = podcastDTO.MapToPodcastDTO(podcast, uid)
	}

	return base.SetData(bookmarksResponse)
}

func (s *UserService) ToggleBookmarkPodcast(userID, podcastID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	pid, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("Invalid Podcast ID", err)
	}

	exists, err := s.BookmarksRepo.IsBookmarked(uid, pid)
	if err != nil {
		return base.SetErrorMessage("Failed to check bookmark", err)
	}

	var action string
	if exists {
		err = s.BookmarksRepo.RemoveBookmark(uid, pid)
		action = "removed"
	} else {
		err = s.BookmarksRepo.AddBookmark(uid, pid)
		action = "added"
	}
	if err != nil {
		return base.SetErrorMessage("Failed to toggle bookmark", err)
	}

	return base.SetSuccessMessage("Bookmark " + action + " successfully")
}

func (s *UserService) GetDownloadedPodcasts(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("Invalid User ID", err)
	}

	downloads, err := s.UserRepo.FindDownloadedPodcasts(uid)
	if err != nil {
		return base.SetErrorMessage("Failed to fetch downloaded podcasts", err)
	}

	response := make([]interface{}, len(downloads))
	for i, podcast := range downloads {
		response[i] = podcastDTO.MapToPodcastDTO(podcast, uid)
	}

	return base.SetData(response)
}
