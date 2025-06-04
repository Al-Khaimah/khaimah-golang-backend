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

// @Summary     Create a new user
// @Description Registers a user with first name, last name, email, categories, and password
// @Tags        users
// @Accept      json
// @Param       signupDTO  body      userDTO.SignupRequestDTO  true  "Signup payload"
// @Success     200        {object}  userDTO.SignupRequestDTO
// @Failure     400        {object}  userDTO.SignupRequestDTO
// @Router /auth/signup [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var signupDTO userDTO.SignupRequestDTO
	if res, ok := base.BindAndValidate(c, &signupDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.CreateUser(&signupDTO)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Login a user
// @Description Authenticates a user and returns a JWT token
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       loginDTO  body      userDTO.LoginRequestDTO  true  "Login payload"
// @Success     200       {object}  userDTO.LoginResponseDTO
// @Failure     400       {object}  echo.HTTPError
// @Router      /auth/login [post]
func (h *UserHandler) LoginUser(c echo.Context) error {
	var loginDTO userDTO.LoginRequestDTO
	if res, ok := base.BindAndValidate(c, &loginDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.LoginUser(&loginDTO)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Logout current user
// @Description Invalidates the current JWT token (if any) and logs out
// @Tags        users
// @Produce     json
// @Success     200  {string}  string  "Logged out successfully"
// @Failure     400  {object}  echo.HTTPError
// @Router      /auth/logout [post]
func (h *UserHandler) LogoutUser(c echo.Context) error {
	response := h.UserService.LogoutUser(c)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Get profile of current user
// @Description Returns profile info for the authenticated user
// @Tags        users
// @Produce     json
// @Success     200  {object}  userDTO.UserProfileDTO
// @Failure     400  {object}  echo.HTTPError
// @Router      /user/profile [get]
func (h *UserHandler) GetUserProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}
	response := h.UserService.GetUserProfile(userID)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Update profile of current user
// @Description Updates first name and/or last name of the authenticated user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       updateDTO  body      userDTO.UpdateProfileDTO  true  "Fields to update"
// @Success     200        {object}  userDTO.UserProfileDTO
// @Failure     400        {object}  echo.HTTPError
// @Router      /user/profile [put]
func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}
	var updateDTO userDTO.UpdateProfileDTO
	if res, ok := base.BindAndValidate(c, &updateDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.UpdateUserProfile(userID, updateDTO)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Update user preferences
// @Description Updates the authenticated user's categories/preferences
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       preferencesDTO  body      userDTO.UpdatePreferencesDTO  true  "New categories list"
// @Success     200             {object}  userDTO.UpdatePreferencesResponseDTO
// @Failure     400             {object}  echo.HTTPError
// @Router      /user/profile/preferences [patch]
func (h *UserHandler) UpdateUserPreferences(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}

	var preferencesDTO userDTO.UpdatePreferencesDTO
	if res, ok := base.BindAndValidate(c, &preferencesDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response, err := h.UserService.UpdateUserPreferences(userID, preferencesDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, base.SetErrorMessage(err.Error()))
	}
	return c.JSON(http.StatusOK, base.SetData(response, "تم تحديث تفضيلات المستخدم بنجاح"))
}

// @Summary     Change password
// @Description Changes the password of the authenticated user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       passwordDTO  body      userDTO.ChangePasswordDTO  true  "Old and new passwords"
// @Success     200          {string}  string  "Password changed successfully"
// @Failure     400          {object}  echo.HTTPError
// @Router      /user/profile/password [patch]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}

	var passwordDTO userDTO.ChangePasswordDTO
	if res, ok := base.BindAndValidate(c, &passwordDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.ChangePassword(userID, passwordDTO)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Delete user's account
// @Description Deletes the authenticated user's account
// @Tags        users
// @Produce     json
// @Success     200  {string}  string  "تم حذف حساب المستخدم بنجاح"
// @Failure     400  {object}  echo.HTTPError
// @Router      /user/profile/delete-my-account [delete]
func (h *UserHandler) DeleteMyAccount(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}

	response := h.UserService.DeleteUser(userID)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Get all users (admin only)
// @Description Returns a list of all users (only accessible by admins)
// @Tags        users
// @Produce     json
// @Success     200  {array}   userDTO.UserProfileDTO
// @Failure     400  {object}  echo.HTTPError
// @Router      /admin/all-users [get]
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	response := h.UserService.GetAllUsers(c)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Mark a user as admin
// @Description Grants admin privileges to a specific user (admin only)
// @Tags        users
// @Produce     json
// @Param       user_id  path      string  true  "ID of the user to promote"
// @Success     200      {string}  string  "User marked as admin"
// @Failure     400      {object}  echo.HTTPError
// @Router      /admin/mark-user-admin/:user_id [post]
func (h *UserHandler) MarkUserAsAdmin(c echo.Context) error {
	userID := c.Param("user_id")

	response := h.UserService.MarkUserAdmin(userID)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Delete a user (admin only)
// @Description Deletes any user by ID (admin only)
// @Tags        users
// @Produce     json
// @Param       id   path      string  true  "ID of the user to delete"
// @Success     200  {string}  string  "User deleted successfully"
// @Failure     400  {object}  echo.HTTPError
// @Router      /admin/user/:id [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID := c.Param("id")

	response := h.UserService.DeleteUser(userID)
	return c.JSON(response.HTTPStatus, response)
}

// @Summary     Get bookmarked podcasts of current user
// @Description Retrieves a list of podcasts bookmarked by the authenticated user
// @Tags        users
// @Produce     json
// @Success     200  {object}  userDTO.GetUserBookmarksResponseDTO
// @Failure     400  {object}  echo.HTTPError
// @Router      /user/bookmarks [get]
func (h *UserHandler) GetUserBookmarks(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}

	response := h.UserService.GetUserBookmarks(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) GetDownloadedPodcasts(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}

	response := h.UserService.GetDownloadedPodcasts(userID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) ToggleBookmarkPodcast(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}
	podcastID := c.Param("podcast_id")

	response := h.UserService.ToggleBookmarkPodcast(userID, podcastID)
	return c.JSON(response.HTTPStatus, response)
}

func (h *UserHandler) CreateSSOUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, base.SetErrorMessage("غير مصرح به"))
	}

	var createSSOUserRequestDTO userDTO.CreateSSOUserRequestDTO
	if res, ok := base.BindAndValidate(c, &createSSOUserRequestDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.UserService.CreateSSOUser(userID, createSSOUserRequestDTO)
	return c.JSON(response.HTTPStatus, response)
}
