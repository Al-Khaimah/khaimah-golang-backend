package middlewares

import (
	"net/http"
	"strings"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	repos "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(authRepo *repos.AuthRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "No token provided"))
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			jwtSecret := config.GetEnv("JWT_SECRET", "alkhaimah123")
			if jwtSecret == "" {
				return c.JSON(http.StatusInternalServerError, base.SetErrorMessage("Server Error", "JWT secret is missing"))
			}

			userID, err := base.ExtractUserIDFromToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid token"))
			}

			authRecord, err := authRepo.FindAuthByUserID(userID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "User authentication not found"))
			}
			if !authRecord.IsActive {
				return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "User is logged out"))
			}

			c.Set("user_id", userID.String())
			return next(c)
		}
	}
}
