package users

import (
	"fmt"

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

func NewUserService(userRepo *repos.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) CreateUser(user *userDTO.SignupRequestDTO) (*userDTO.SignupResponseDTO, error) {
	existingUser, err := s.UserRepo.FindOneByEmail(user.Email)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email: %s exists", user.Email)
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
		return nil, fmt.Errorf("service error: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUserAuth := &models.IamAuth{
		UserID:   createdUser.ID,
		Password: string(hashedPassword),
	}

	err = s.AuthRepo.CreateUserAuth(newUserAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth for user ID %s: %w", createdUser.ID.String(), err)
	}

	// return
}

func convertCategories(categoryIDs []string) []categories.Category {
	categoryList := make([]categories.Category, len(categoryIDs))
	for i, id := range categoryIDs {
		var category categories.Category
		category.ID = uuid.MustParse(id)
		categoryList[i] = category
	}
	return categoryList
}
