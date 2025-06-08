package users

import (
	"fmt"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base/utils"
	userDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/enums"
	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	repos "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	podcastDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/dtos"
)

type UserService struct {
	UserRepo      *repos.UserRepository
	AuthRepo      *repos.AuthRepository
	BookmarksRepo *repos.BookmarkRepository
}

func NewUserService(userRepo *repos.UserRepository, authRepo *repos.AuthRepository, bookmarksRepo *repos.BookmarkRepository) *UserService {
	return &UserService{
		UserRepo:      userRepo,
		AuthRepo:      authRepo,
		BookmarksRepo: bookmarksRepo,
	}
}

func (s *UserService) CreateUser(user *userDTO.SignupRequestDTO) base.Response {
	existingUser, err := s.UserRepo.FindOneByEmail(user.Email)
	if err != nil {
		return base.SetErrorMessage("خطأ في قاعدة البيانات")
	}
	if existingUser != nil {
		return base.SetErrorMessage("هذا البريد الإلكتروني قيد الاستخدام بالفعل")
	}

	categories := utils.ConvertIDsToCategories(user.Categories)
	newUser := &models.User{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      utils.FormatEmail(user.Email),
		Categories: categories,
	}

	createdUser, err := s.UserRepo.CreateUser(newUser)
	if err != nil {
		return base.SetErrorMessage("فشل في إنشاء المستخدم")
	}

	if createdUser == nil {
		return base.SetErrorMessage("فشل في إنشاء المستخدم")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return base.SetErrorMessage("فشل في تشفير كلمة المرور")
	}

	newUserAuth := &models.IamAuth{
		UserID:   createdUser.ID,
		Password: string(hashedPassword),
		IsActive: true,
	}

	if err := s.AuthRepo.CreateUserAuth(newUserAuth); err != nil {
		return base.SetErrorMessage("فشل في إنشاء توثيق المستخدم")
	}

	token, err := GenerateJWT(createdUser)
	if err != nil {
		return base.SetErrorMessage("فشل في إنشاء الرمز")
	}
	userResponse := userDTO.SignupResponseDTO{
		ID:         createdUser.ID.String(),
		FirstName:  createdUser.FirstName,
		LastName:   createdUser.LastName,
		Email:      createdUser.Email,
		Categories: createdUser.Categories,
		Token:      token,
		ExpiresAt:  "never",
	}

	slackMessage := fmt.Sprintf("🚀 New user account created:\n%s (%s)", createdUser.FirstName, createdUser.Email)
	_ = utils.SendSlackNotification(slackMessage)

	return base.SetData(userResponse, "تم انشاء الحساب بنجاح")
}

func (s *UserService) LoginUser(user *userDTO.LoginRequestDTO) base.Response {
	existingUser, err := s.UserRepo.FindOneByEmail(user.Email)
	if err != nil {
		return base.SetErrorMessage("خطأ في قاعدة البيانات")
	}
	if existingUser == nil {
		return base.SetErrorMessage("البريد الإلكتروني غير موجود")
	}

	userAuth, err := s.AuthRepo.FindAuthByUserID(existingUser.ID)
	if err != nil {
		return base.SetErrorMessage("خطأ في التوثيق")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(user.Password)); err != nil {
		return base.SetErrorMessage("كلمة المرور غير صحيحة")
	}

	token, err := GenerateJWT(existingUser)
	if err != nil {
		return base.SetErrorMessage("فشل في إنشاء الرمز")
	}

	userAuth.IsActive = true
	if err := s.AuthRepo.UpdateAuth(userAuth); err != nil {
		return base.SetErrorMessage("فشل في تحديث التوثيق")
	}

	loginResponse := userDTO.LoginResponseDTO{
		ID:         existingUser.ID.String(),
		FirstName:  existingUser.FirstName,
		LastName:   existingUser.LastName,
		Email:      existingUser.Email,
		Categories: existingUser.Categories,
		Token:      token,
		ExpiresAt:  "never",
	}

	return base.SetData(loginResponse, "Logged in successfully")
}

func GenerateJWT(user *models.User) (string, error) {
	jwtSecret := config.GetEnv("JWT_SECRET", "alkhaimah123")
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   utils.FormatEmail(user.Email),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (s *UserService) LogoutUser(c echo.Context) base.Response {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return base.SetErrorMessage("غير مصرح به")
	}

	userID, err := utils.ExtractUserIDFromToken(token)
	if err != nil {
		return base.SetErrorMessage("رمز غير صالح")
	}

	authRecord, err := s.AuthRepo.FindAuthByUserID(userID)
	if err != nil {
		return base.SetErrorMessage("فشل في العثور على توثيق المستخدم")
	}

	authRecord.IsActive = false
	err = s.AuthRepo.UpdateAuth(authRecord)
	if err != nil {
		return base.SetErrorMessage("فشل في تسجيل الخروج")
	}

	return base.SetSuccessMessage("Successfully logged out")
}

func (s *UserService) GetUserProfile(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("فشل في استرجاع ملف تعريف المستخدم")
	}
	if user == nil {
		return base.SetErrorMessage("لم يتم العثور على مستخدم بهذا الرقم التعريفي")
	}

	profileResponse := userDTO.UserProfileDTO{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Categories: user.Categories,
	}

	return base.SetData(profileResponse)
}

func (s *UserService) UpdateUserProfile(userID string, updateData userDTO.UpdateProfileDTO) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("خطأ في قاعدة البيانات")
	}
	if user == nil {
		return base.SetErrorMessage("لم يتم العثور على مستخدم بهذا الرقم التعريفي")
	}

	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}

	err = s.UserRepo.UpdateUser(user)
	if err != nil {
		return base.SetErrorMessage("فشل في تحديث الملف الشخصي")
	}

	profileResponse := userDTO.UserProfileDTO{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return base.SetData(profileResponse, "Profile updated successfully")
}

func (s *UserService) UpdateUserPreferences(userID string, updateData userDTO.UpdatePreferencesDTO) (*userDTO.UpdatePreferencesDTO, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("الرقم التعريفي للمستخدم غير صالح")
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return nil, fmt.Errorf("خطأ في قاعدة البيانات")
	}
	if user == nil {
		return nil, fmt.Errorf("لم يتم العثور على مستخدم بهذا الرقم التعريفي")
	}

	newCategories := utils.ConvertIDsToCategories(updateData.Categories)

	err = s.UserRepo.UpdateUserPreferences(user, newCategories)
	if err != nil {
		return nil, fmt.Errorf("فشل في تحديث التفضيلات")
	}

	preferencesResponse := userDTO.UpdatePreferencesDTO{
		Categories: utils.ConvertCategoriesToStringIDs(user.Categories),
	}
	return &preferencesResponse, nil
}

func (s *UserService) ChangePassword(userID string, req userDTO.ChangePasswordDTO) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	userAuth, err := s.AuthRepo.FindAuthByUserID(uid)
	if err != nil {
		return base.SetErrorMessage("خطأ في قاعدة البيانات")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAuth.Password), []byte(req.OldPassword))
	if err != nil {
		return base.SetErrorMessage("كلمة المرور القديمة غير صحيحة")
	}

	if req.OldPassword == req.NewPassword {
		return base.SetErrorMessage("يجب أن تكون كلمة المرور الجديدة مختلفة عن كلمة المرور القديمة")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return base.SetErrorMessage("فشل في تشفير كلمة المرور")
	}

	userAuth.Password = string(hashedPassword)

	err = s.AuthRepo.UpdateAuth(userAuth)
	if err != nil {
		return base.SetErrorMessage("فشل في تغيير كلمة المرور")
	}

	return base.SetSuccessMessage("Password changed successfully")
}

func (s *UserService) GetAllUsers(c echo.Context) base.Response {
	users, err := s.UserRepo.FindAllUsers()
	if err != nil {
		return base.SetErrorMessage("فشل في استرجاع المستخدمين")
	}

	var userResponses []interface{}
	for _, user := range users {
		userResponses = append(userResponses, userDTO.UserProfileDTO{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}
	return base.SetDataPaginated(c, userResponses)
}

func (s *UserService) MarkUserAdmin(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("فشل في استرجاع المستخدم")
	}
	if user == nil {
		return base.SetErrorMessage("المستخدم غير موجود")
	}

	user.UserType = users.UserTypeAdmin
	if err := s.UserRepo.UpdateUser(user); err != nil {
		return base.SetErrorMessage("فشل في تحديث المستخدم")
	}

	return base.SetSuccessMessage("User marked as admin successfully")
}

func (s *UserService) DeleteUser(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("خطأ في قاعدة البيانات")
	}
	if user == nil {
		return base.SetErrorMessage("لم يتم العثور على مستخدم بهذا الرقم التعريفي")
	}

	if user.UserType == users.UserTypeAdmin {
		return base.SetErrorMessage("لا يمكنك حذف حساب مستخدم يمتلك صلاحيات المشرف")
	}

	if user.DeletedAt.Valid {
		return base.SetErrorMessage("تم حذف هذا الحساب مسبقاً")
	}

	err = s.UserRepo.DeleteUser(uid)
	if err != nil {
		return base.SetErrorMessage("فشل في حذف المستخدم")
	}

	return base.SetSuccessMessage("تم حذف حساب المستخدم بنجاح")
}

func (s *UserService) GetUserCategoriesIDs(userID string) ([]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	userCategories, err := s.UserRepo.FindUserCategories(uid)
	if err != nil {
		return nil, err
	}
	if userCategories == nil {
		return nil, nil
	}

	categoryIDs := make([]string, len(userCategories))
	for i, userCategories := range userCategories {
		categoryIDs[i] = userCategories.ID.String()
	}

	return categoryIDs, nil
}

func (s *UserService) GetUserBookmarks(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	bookmarks, err := s.BookmarksRepo.FindUserBookmarks(uid)
	if err != nil {
		return base.SetErrorMessage("فشل في استرجاع الإشارات المرجعية")
	}

	bookmarksResponse := make([]interface{}, len(bookmarks))
	for i, podcast := range bookmarks {
		bookmarksResponse[i] = podcastDTO.MapToPodcastDTO(podcast, uid)
	}

	return base.SetData(bookmarksResponse)
}

func (s *UserService) ToggleBookmarkPodcast(userID, podcastID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	pid, err := uuid.Parse(podcastID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للبودكاست غير صالح")
	}

	exists, err := s.BookmarksRepo.IsBookmarked(uid, pid)
	if err != nil {
		return base.SetErrorMessage("فشل في التحقق من الإشارة المرجعية")
	}

	var action string
	if exists {
		err = s.BookmarksRepo.RemoveBookmark(uid, pid)
		action = "removed"
	} else {
		err = s.BookmarksRepo.AddBookmark(uid, pid)
		action = "added"
	}
	if err != nil {
		return base.SetErrorMessage("فشل في تبديل الإشارة المرجعية")
	}

	return base.SetSuccessMessage("Bookmark " + action + " successfully")
}

func (s *UserService) GetDownloadedPodcasts(userID string) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	downloads, err := s.UserRepo.FindDownloadedPodcasts(uid)
	if err != nil {
		return base.SetErrorMessage("فشل في استرجاع البودكاست التي تم تنزيلها")
	}

	response := make([]interface{}, len(downloads))
	for i, podcast := range downloads {
		response[i] = podcastDTO.MapToPodcastDTO(podcast, uid)
	}

	return base.SetData(response)
}

func (s *UserService) CreateSSOUser(userID string, createSSOUserRequestDTO userDTO.CreateSSOUserRequestDTO) base.Response {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return base.SetErrorMessage("الرقم التعريفي للمستخدم غير صالح")
	}

	user, err := s.UserRepo.FindOneByID(uid)
	if err != nil {
		return base.SetErrorMessage("خطأ في قاعدة البيانات")
	}

	user.FirstName = createSSOUserRequestDTO.FirstName

	err = s.UserRepo.UpdateUser(user)
	if err != nil {
		return base.SetErrorMessage("فشل في تحديث الملف الشخصي")
	}

	updatePreferencesDTO := userDTO.UpdatePreferencesDTO{
		Categories: createSSOUserRequestDTO.Categories,
	}

	updatedPreferences, err := s.UpdateUserPreferences(userID, updatePreferencesDTO)
	if err != nil {
		return base.SetErrorMessage("فشل في تحديث التفضيلات")
	}

	userData := userDTO.CreateSSOUserResponseDTO{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		Email:      user.Email,
		Categories: utils.ConvertIDsToCategories(updatedPreferences.Categories),
	}

	return base.SetData(userData, "تم انشاء حساب المستخدم بنجاح")
}
