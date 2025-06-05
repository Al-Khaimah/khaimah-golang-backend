package users

import (
	"net/http"

	base "github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	authDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	AuthService *users.AuthService
}

func NewAuthHandler(authService *users.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) OAuthLogin(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, base.SetErrorMessage("Unauthorized", "No param token provided"))
	}

	var oAuthRequestDTO authDTO.OAuthRequestDTO
	if res, ok := base.BindAndValidate(c, &oAuthRequestDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.AuthService.Login(c.Request().Context(), oAuthRequestDTO.Provider, token)
	if response.Errors != nil {
		code := http.StatusUnauthorized
		if response.Errors == "unsupported provider" {
			code = http.StatusBadRequest
		}
		return c.JSON(code, response)
	}

	return c.JSON(http.StatusOK, response)
}
