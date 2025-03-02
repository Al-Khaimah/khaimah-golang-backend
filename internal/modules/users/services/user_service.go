package users

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	userDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	repos "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
