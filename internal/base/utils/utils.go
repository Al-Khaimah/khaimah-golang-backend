package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
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

func FormatEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ConvertCategoriesToString(categories []models.Category) []string {
	categoryIDs := make([]string, len(categories))
	for i, category := range categories {
		categoryIDs[i] = category.ID.String()
	}
	return categoryIDs
}

func ConvertIDsToCategories(categoryIDs []string) []models.Category {
	var categoryList []models.Category

	for _, id := range categoryIDs {
		uid, err := uuid.Parse(id)
		if err != nil {
			continue
		}

		categoryRepo := categories.NewCategoryRepository(config.GetDB())
		category, err := categoryRepo.FindCategoryByID(uid)
		if err != nil || category == nil {
			continue
		}

		categoryList = append(categoryList, *category)
	}

	return categoryList
}
