package users

import (
	"net/http"

	base "github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	userDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	userService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/services"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService *userService.UserService
}

func NewUserHandler(userService *userService.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var signupDTO userDTO.SignupRequestDTO
	if res, ok := base.BindAndValidate(c, &signupDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.CreateUser(&signupDTO)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) LoginUser(c echo.Context) error {
	var loginDTO userDTO.LoginRequestDTO
	if res, ok := base.BindAndValidate(c, &loginDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.LoginUser(&loginDTO)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) LogoutUser(c echo.Context) error {
	response := h.UserService.LogoutUser(c)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) GetUserProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}
	response := h.UserService.GetUserProfile(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}
	var updateDTO userDTO.UpdateProfileDTO
	if res, ok := base.BindAndValidate(c, &updateDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.UpdateUserProfile(userID, updateDTO)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) UpdateUserPreferences(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var preferencesDTO userDTO.UpdatePreferencesDTO
	if res, ok := base.BindAndValidate(c, &preferencesDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.UpdateUserPreferences(userID, preferencesDTO)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}
	var passwordDTO userDTO.ChangePasswordDTO
	if res, ok := base.BindAndValidate(c, &passwordDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.ChangePassword(userID, passwordDTO)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) DeleteMyAccount(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("Unauthorized", "Invalid or missing user ID"))
	}

	response := h.UserService.DeleteUser(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	response := h.UserService.GetAllUsers(c)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) MarkUserAsAdmin(c echo.Context) error {
	userID := c.Param("user_id")

	response := h.UserService.MarkUserAdmin(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID := c.Param("id")

	response := h.UserService.DeleteUser(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) GetUserBookmarks(c echo.Context) error {
	userID := c.Get("user_id").(string)

	response := h.UserService.GetUserBookmarks(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) GetDownloadedPodcasts(c echo.Context) error {
	userID := c.Get("user_id").(string)

	response := h.UserService.GetDownloadedPodcasts(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) ToggleBookmarkPodcast(c echo.Context) error {
	userID := c.Get("user_id").(string)
	podcastID := c.Param("podcast_id")

	response := h.UserService.ToggleBookmarkPodcast(userID, podcastID)
	return c.JSON(response.HTTPStatus, response)
}
