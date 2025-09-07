package users

import (
	"context"
	"fmt"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base/utils"
	userDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	repos "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
)

// OTPService handles OTP-related operations
type OTPService struct {
	UserRepo  *repos.UserRepository
	AuthRepo  *repos.AuthRepository
	JWTSecret []byte
}

// NewOTPService creates a new OTP service
func NewOTPService(userRepo *repos.UserRepository, authRepo *repos.AuthRepository, jwtSecret []byte) *OTPService {
	return &OTPService{
		UserRepo:  userRepo,
		AuthRepo:  authRepo,
		JWTSecret: jwtSecret,
	}
}

// SendOTP sends an OTP to the provided email or mobile number
func (s *OTPService) SendOTP(req *userDTO.SendOTPRequestDTO) base.Response {
	ctx := context.Background()
	if req.Email == "" && req.Mobile == "" {
		return base.SetErrorMessage("email or mobile is required")
	}

	var identifier string
	var existingUser *models.User
	var err error

	if req.Email != "" {
		req.Email = utils.FormatEmail(req.Email)
		identifier = req.Email
		existingUser, err = s.UserRepo.FindOneByEmail(req.Email)
	} else {
		req.Mobile = utils.FormatMobileNumber(req.Mobile)
		identifier = req.Mobile
		existingUser, err = s.UserRepo.FindOneByMobile(req.Mobile)
	}

	if err != nil {
		return base.SetErrorMessage("خطأ في قاعدة البيانات")
	}

	otp := utils.GenerateOTP()

	if err := utils.StoreOTP(ctx, identifier, otp); err != nil {
		return base.SetErrorMessage("فشل في تخزين رمز التحقق")
	}

	if existingUser == nil {
		if req.FirstName == "" {
			return base.SetErrorMessage("first_name is required for new users")
		}
		if len(req.Categories) == 0 {
			return base.SetErrorMessage("categories are required for new users")
		}

		categories := utils.ConvertIDsToCategories(req.Categories)
		newUser := &models.User{
			FirstName:  req.FirstName,
			Categories: categories,
		}

		if req.Email != "" {
			newUser.Mobile = ""
			newUser.Email = req.Email
		} else if req.Mobile != "" {
			newUser.Email = fmt.Sprintf("mobile_%s@placeholder.com", req.Mobile)
			newUser.Mobile = req.Mobile
		}

		createdUser, err := s.UserRepo.CreateUser(newUser)
		if err != nil {
			return base.SetErrorMessage("فشل في إنشاء المستخدم", err)
		}

		newUserAuth := &models.IamAuth{
			UserID:   createdUser.ID,
			IsActive: false,
		}

		if err := s.AuthRepo.CreateUserAuth(newUserAuth); err != nil {
			return base.SetErrorMessage("فشل في إنشاء توثيق المستخدم")
		}
	}

	var sendErr error
	if req.Email != "" {
		sendErr = utils.SendEmailOTP(req.Email, otp, req.FirstName)
	} else {
		formattedMobile := utils.FormatMobileNumber(req.Mobile)
		sendErr = utils.SendMobileOTP(formattedMobile, otp, req.FirstName)
	}

	if sendErr != nil {
		return base.SetErrorMessage(fmt.Sprintf("فشل في إرسال رمز التحقق: %v", sendErr))
	}

	return base.SetSuccessMessage("تم إرسال رمز التحقق بنجاح")
}

// VerifyOTP verifies the OTP and returns a JWT token if valid
func (s *OTPService) VerifyOTP(req *userDTO.VerifyOTPRequestDTO) base.Response {
	ctx := context.Background()
	// Validate that either email or mobile is provided
	if req.Email == "" && req.Mobile == "" {
		return base.SetErrorMessage("email or mobile is required")
	}

	var identifier string
	var user *models.User
	var err error

	if req.Email != "" {
		req.Email = utils.FormatEmail(req.Email)
		identifier = req.Email
		user, err = s.UserRepo.FindOneByEmail(req.Email)
	} else {
		formattedMobile := utils.FormatMobileNumber(req.Mobile)
		identifier = formattedMobile
		user, err = s.UserRepo.FindOneByMobile(formattedMobile)
	}

	if err != nil {
		return base.SetErrorMessage("خطأ في قاعدة البيانات")
	}

	if user == nil {
		return base.SetErrorMessage("المستخدم غير موجود")
	}

	isValid, err := utils.VerifyOTP(ctx, identifier, req.OTP)
	if err != nil {
		return base.SetErrorMessage(fmt.Sprintf("فشل في التحقق من الرمز: %v", err))
	}

	if !isValid {
		return base.SetErrorMessage("رمز التحقق غير صالح أو منتهي الصلاحية")
	}

	_ = utils.DeleteOTP(ctx, identifier)

	userAuth, err := s.AuthRepo.FindAuthByUserID(user.ID)
	if err != nil {
		return base.SetErrorMessage("خطأ في التوثيق")
	}

	userAuth.IsActive = true
	if err := s.AuthRepo.UpdateAuth(userAuth); err != nil {
		return base.SetErrorMessage("فشل في تحديث حالة المستخدم")
	}

	tokenString, err := GenerateJWT(user)
	if err != nil {
		return base.SetErrorMessage("فشل في إنشاء الرمز")
	}

	response := userDTO.VerifyOTPResponseDTO{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		Email:      user.Email,
		Mobile:     user.Mobile,
		Categories: user.Categories,
		Token:      tokenString,
		ExpiresAt:  "never",
	}

	return base.SetData(response, "تم تسجيل الدخول بنجاح")
}
