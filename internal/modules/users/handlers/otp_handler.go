package users

import (
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	userDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	userService "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/services"
	"github.com/labstack/echo/v4"
)

type OTPHandler struct {
	otpService *userService.OTPService
}

func NewOTPHandler(otpService *userService.OTPService) *OTPHandler {
	return &OTPHandler{
		otpService: otpService,
	}
}

// SendOTP godoc
// @Summary Send OTP to email or mobile
// @Description Send a one-time password to the provided email or mobile number
// @Tags auth
// @Accept json
// @Produce json
// @Param request body userDTO.SendOTPRequestDTO true "Send OTP request"
// @Success 200 {object} base.Response
// @Failure 400 {object} base.Response
// @Router /auth/send-otp [post]
func (h *OTPHandler) SendOTP(c echo.Context) error {
	var req userDTO.SendOTPRequestDTO
	if res, ok := base.BindAndValidate(c, &req); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.otpService.SendOTP(&req)
	return c.JSON(response.HTTPStatus, response)
}

// VerifyOTP godoc
// @Summary Verify OTP and get token
// @Description Verify the OTP sent to email or mobile and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body userDTO.VerifyOTPRequestDTO true "Verify OTP request"
// @Success 200 {object} base.Response
// @Failure 400 {object} base.Response
// @Router /auth/verify-otp [post]
func (h *OTPHandler) VerifyOTP(c echo.Context) error {
	var req userDTO.VerifyOTPRequestDTO
	if res, ok := base.BindAndValidate(c, &req); !ok {
		return c.JSON(res.HTTPStatus, res)
	}

	response := h.otpService.VerifyOTP(&req)
	return c.JSON(response.HTTPStatus, response)
}
