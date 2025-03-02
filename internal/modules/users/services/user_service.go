package users

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	userDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	repos "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
)

var jwtSecret = config.GetEnv("JWT_SECRET", "alkhaimah123")

type UserService struct {
	UserRepo *repos.UserRepository
	AuthRepo *repos.AuthRepository
}

func NewUserService(userRepo *repos.UserRepository, authRepo *repos.AuthRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		AuthRepo: authRepo,
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

	categories := convertCategories(user.Categories)
	newUser := &models.User{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
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
	}

	if err := s.AuthRepo.CreateUserAuth(newUserAuth); err != nil {
		return base.SetErrorMessage("Failed to create user authentication", err)
	}

	userResponse := userDTO.SignupResponseDTO{
		ID:        createdUser.ID.String(),
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Email:     createdUser.Email,
	}

	return base.SetData(userResponse)
}

func convertCategories(categoryIDs []string) []categories.Category {
	if categoryIDs == nil || len(categoryIDs) == 0 {
		return []categories.Category{}
	}

	categoryList := make([]categories.Category, len(categoryIDs))
	for i, id := range categoryIDs {
		var category categories.Category
		category.ID = uuid.MustParse(id)
		categoryList[i] = category
	}
	return categoryList
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
		ID:        existingUser.ID.String(),
		FirstName: existingUser.FirstName,
		LastName:  existingUser.LastName,
		Email:     existingUser.Email,
		Token:     token,
		ExpiresAt: "never",
	}

	return base.SetData(loginResponse)
}

func generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		/*"exp":   time.Now().Add(time.Hour * 24).Unix(),*/ //if we wanted to set ExpiresAt
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *UserService) LogoutUser(c echo.Context) base.Response {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return base.SetErrorMessage("Unauthorized", "No token provided")
	}

	userID, err := s.AuthRepo.ExtractUserIDFromToken(token)
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
