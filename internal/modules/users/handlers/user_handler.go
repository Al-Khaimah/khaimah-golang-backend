package users

import (
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

	userResponse := h.UserService.CreateUser(&signupDTO)
	return c.JSON(userResponse.HTTPStatus, userResponse)
}

func (h *UserHandler) LoginUser(c echo.Context) error {
	var loginDTO userDTO.LoginRequestDTO
	if res, ok := base.BindAndValidate(c, &loginDTO); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	loginResponse := h.UserService.LoginUser(&loginDTO)
	return c.JSON(loginResponse.HTTPStatus, loginResponse)
}

func (h *UserHandler) LogoutUser(c echo.Context) error {
	userResponse := h.UserService.LogoutUser(c)
	return c.JSON(userResponse.HTTPStatus, userResponse)
}
