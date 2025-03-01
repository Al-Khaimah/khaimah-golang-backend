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
	if err := base.BindAndValidate(c, &signupDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.UserService.CreateUser(&signupDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}
