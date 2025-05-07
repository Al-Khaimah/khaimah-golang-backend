package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func SendSlackNotification(message string) error {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		return fmt.Errorf("SLACK_WEBHOOK_URL environment variable not set")
	}

	slackPayload := map[string]string{"text": message}
	payloadBytes, err := json.Marshal(slackPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create Slack HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Slack notification request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		bodyString := ""
		if readErr == nil {
			bodyString = string(bodyBytes)
		}
		return fmt.Errorf("failed to send Slack notification: received status code %d, response: %s", resp.StatusCode, bodyString)
	}

	io.Copy(io.Discard, resp.Body)

	return nil
}
