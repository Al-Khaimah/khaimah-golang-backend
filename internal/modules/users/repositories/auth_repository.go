package users

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"strings"

	"gorm.io/gorm"

	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) CreateUserAuth(userAuth *models.IamAuth) error {
	result := r.DB.Create(userAuth)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

func (r *AuthRepository) FindAuthByUserID(userID uuid.UUID) (*models.IamAuth, error) {
	var auth models.IamAuth
	result := r.DB.Where("user_id = ?", userID).First(&auth)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find authentication record: %w", result.Error)
	}
	return &auth, nil
}

func (r *AuthRepository) UpdateAuth(auth *models.IamAuth) error {
	result := r.DB.Save(auth)
	if result.Error != nil {
		return fmt.Errorf("failed to update authentication record: %w", result.Error)
	}
	return nil
}

func (r *AuthRepository) ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID in token")
	}

	return userID, nil
}
