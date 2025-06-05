package users

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	resend "github.com/resend/resend-go/v2"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	authDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"
	models "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	userRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"google.golang.org/api/idtoken"
)

type AuthService struct {
	userRepo  *userRepository.UserRepository
	authRepo  *userRepository.AuthRepository
	providers map[string]Provider
	jwtSecret []byte
}

type Provider interface {
	Exchange(ctx context.Context, idToken string) (string, error)
	GetClientSecret(ctx context.Context) (string, error)
}

type GoogleProvider struct {
	ClientID string
}

type AppleProvider struct {
	JWKSUrl    string
	ClientID   string
	TeamID     string
	KeyID      string
	PrivateKey string
}

func NewAuthService(repo *userRepository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: repo,
		providers: map[string]Provider{
			"google": NewGoogleProvider(),
			"apple":  NewAppleProvider(),
		},
		jwtSecret: []byte(os.Getenv("JWT_SECRET")),
	}
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{ClientID: os.Getenv("GOOGLE_CLIENT_ID")}
}

func NewAppleProvider() *AppleProvider {
	return &AppleProvider{
		JWKSUrl:    "https://appleid.apple.com/auth/keys",
		ClientID:   os.Getenv("APPLE_CLIENT_ID"),
		TeamID:     os.Getenv("APPLE_TEAM_ID"),
		KeyID:      os.Getenv("APPLE_KEY_ID"),
		PrivateKey: os.Getenv("APPLE_PRIVATE_KEY"),
	}
}

func (g *GoogleProvider) Exchange(ctx context.Context, token string) (string, error) {
	p, err := idtoken.Validate(ctx, token, g.ClientID)
	if err != nil {
		return "", fmt.Errorf("google token invalid: %w", err)
	}
	email, _ := p.Claims["email"].(string)
	return email, nil
}

func (a *AppleProvider) Exchange(ctx context.Context, token string) (string, error) {
	set, err := jwk.Fetch(ctx, a.JWKSUrl)
	if err != nil {
		return "", fmt.Errorf("apple jwk fetch failed: %w", err)
	}

	tok, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid")
		}
		key, found := set.LookupKeyID(kid)
		if !found {
			return nil, fmt.Errorf("jwk %s not found", kid)
		}
		var raw interface{}
		if err := key.Raw(&raw); err != nil {
			return nil, fmt.Errorf("invalid key: %w", err)
		}
		return raw, nil
	})

	if err != nil || !tok.Valid {
		return "", fmt.Errorf("apple token invalid: %w", err)
	}

	claims := tok.Claims.(jwt.MapClaims)

	if a.ClientID != "" {
		aud, ok := claims["aud"].(string)
		if !ok || aud != a.ClientID {
			return "", fmt.Errorf("invalid audience in token: %v", claims["aud"])
		}
	}

	iss, ok := claims["iss"].(string)
	if !ok || iss != "https://appleid.apple.com" {
		return "", fmt.Errorf("invalid issuer in token: %v", claims["iss"])
	}

	email, _ := claims["email"].(string)
	if email == "" {
		email = fmt.Sprint(claims["sub"])
	}

	return email, nil
}

func (g *GoogleProvider) GetClientSecret(ctx context.Context) (string, error) {
	return "", nil
}

func (a *AppleProvider) GetClientSecret(ctx context.Context) (string, error) {
	if a.PrivateKey == "" || a.KeyID == "" || a.TeamID == "" || a.ClientID == "" {
		return "", fmt.Errorf("missing required Apple credentials")
	}

	block, _ := pem.Decode([]byte(a.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": a.TeamID,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour * 24).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": a.ClientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = a.KeyID

	clientSecret, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign client secret: %w", err)
	}

	return clientSecret, nil
}

func (s *AuthService) Login(ctx context.Context, providerKey, token string) base.Response {
	prov, ok := s.providers[providerKey]
	if !ok {
		return base.SetErrorMessage("unsupported provider")
	}

	email, err := prov.Exchange(ctx, token)
	if err != nil {
		return base.SetErrorMessage("invalid token: " + err.Error())
	}

	user, err := s.userRepo.FindOrCreateByEmail(email)
	if err != nil {
		return base.SetErrorMessage("db error")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.jwtSecret)
	if err != nil {
		return base.SetErrorMessage("sign token error")
	}
	return base.SetData(map[string]string{"token": signed}, "login successful")
}

func (s *AuthService) ValidateWithProvider(ctx context.Context, providerKey, token string) (bool, error) {
	prov, ok := s.providers[providerKey]
	if !ok {
		return false, fmt.Errorf("unsupported provider")
	}

	if providerKey == "apple" {
		clientSecret, err := prov.GetClientSecret(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to generate client secret: %w", err)
		}
		_ = clientSecret
	}

	_, err := prov.Exchange(ctx, token)
	return err == nil, err
}

func (s *AuthService) SendOTPViaSMS(sendOTPViaSMSRequestDTO authDTO.SendOTPViaSMSRequestDTO) error {
	user, err := s.userRepo.FindOneByPhoneNumber(sendOTPViaSMSRequestDTO.Phonenumber)
	if err != nil {
		return fmt.Errorf("لم يتم العثور على مستخدم بالرقم الجوال")
	}

	userAuth, err := s.authRepo.FindAuthByUserID(user.ID)
	if err != nil {
		return fmt.Errorf("حدث خطأ اثناء العثور على بيانات المستخدم")
	}

	otp, err := s.issueOTP(userAuth)
	if err != nil {
		return fmt.Errorf("حدث خطأ أثناء إنشاء رمز التحقق")
	}

	phoneNumbers := []string{user.PhoneNumber}
	message := fmt.Sprintf("رمز التحقق الخاص بك %d", otp) // TODO: use appropriate sms message
	err = s.SendSMS(phoneNumbers, message)
	if err != nil {
		return fmt.Errorf("حدث حطأ اثناء ارسال رسالة نصية الى المستخدم")
	}

	return nil
}

func (s *AuthService) SendOTPViaEmail(ctx context.Context, sendOTPViaEmailRequestDTO authDTO.SendOTPViaEmailRequestDTO) error {
	user, err := s.userRepo.FindOneByEmail(sendOTPViaEmailRequestDTO.Email)
	if err != nil {
		return fmt.Errorf("لم يتم العثور على مستخدم بالرقم الجوال")
	}

	userAuth, err := s.authRepo.FindAuthByUserID(user.ID)
	if err != nil {
		return fmt.Errorf("حدث خطأ اثناء العثور على بيانات المستخدم")
	}

	otp, err := s.issueOTP(userAuth)
	if err != nil {
		return fmt.Errorf("حدث خطأ أثناء إنشاء رمز التحقق")
	}

	email := []string{user.Email}
	err = s.SendEmail(ctx, email, otp)
	if err != nil {
		return fmt.Errorf("حدث حطأ اثناء ارسال رسالة نصية الى المستخدم")
	}

	return nil
}

func (s *AuthService) issueOTP(userAuth *models.IamAuth) (int, error) {
	if time.Now().After(userAuth.ExpiresAt) {
		userAuth.FailedAttempts = 0
	}

	if userAuth.FailedAttempts >= 3 && time.Now().Before(userAuth.ExpiresAt) {
		minutesLeft := max(int(time.Until(userAuth.ExpiresAt).Minutes()), 0)
		minuteWord := getMinuteWord(minutesLeft)
		return 0, fmt.Errorf("لقد تجاوزت عدد المحاولات، حاول بعد %d %s", minutesLeft, minuteWord)
	}

	otp, err := generateOTP()
	if err != nil {
		return 0, fmt.Errorf("حدث خطأ أثناء إنشاء رمز التحقق")
	}

	userAuth.OTP = otp
	userAuth.ExpiresAt = time.Now().Add(15 * time.Minute)
	userAuth.FailedAttempts = 0

	if err := s.authRepo.UpdateAuth(userAuth); err != nil {
		return 0, fmt.Errorf("حدث خطأ أثناء حفظ رمز التحقق")
	}

	return otp, nil
}

func (s *AuthService) SendSMS(phoneNumbers []string, message string) error {
	data := authDTO.SendSMSRequest{
		To:      phoneNumbers,
		Message: message,
		Token:   os.Getenv("SMS_PROVIDER_TOKEN"),
	}

	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "https://api.sendmsg.dev/message/batch", bytes.NewBuffer(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func (s *AuthService) SendEmail(ctx context.Context, emails []string, otp int) error {
	client := resend.NewClient(os.Getenv("EMAIL_PROVIDER_TOKEN"))

	htmlContent := fmt.Sprintf(`
		<p>رمز التحقق الخاص بك هو:</p>
		<h2>%d</h2>
		<p>يرجى استخدامه خلال 15 دقيقة.</p>
	`, otp)

	params := &resend.SendEmailRequest{
		From:    os.Getenv("EMAIL_FROM"),
		To:      emails,
		Subject: "رمز التحقق من الخيمة", // TODO: Replace with proper subject
		Html:    htmlContent,
	}

	_, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("حدث حطأ اثناء ارسال بريد الكتروني الى المستخدم")
	}
	return nil
}

func getMinuteWord(n int) string {
	switch {
	case n == 1:
		return "دقيقة"
	case n == 2:
		return "دقيقتين"
	case n >= 3 && n <= 10:
		return "دقائق"
	default:
		return "دقيقة"
	}
}

func generateOTP() (int, error) {
	firstDigit := rand.Intn(9) + 1

	remaining := rand.Intn(1000)

	otp := firstDigit*1000 + remaining
	return otp, nil
}

func (s *AuthService) VerifyOTP(verifyOTPRequestDTO authDTO.VerifyOTPRequestDTO) error {
	// reset failed attempts to zero upon 1- successful otp, 2- last otp is expired
	var err error
	var user *models.User
	if verifyOTPRequestDTO.Identifier == authDTO.IdentifierEmail {
		user, err = s.userRepo.FindOneByEmail(verifyOTPRequestDTO.Type)
	} else {
		user, err = s.userRepo.FindOneByPhoneNumber(verifyOTPRequestDTO.Type)
	}
	if err != nil {
		return fmt.Errorf("لم يتم العثور على مستخدم")
	}

	userAuth, err := s.authRepo.FindAuthByUserID(user.ID)
	if err != nil {
		return fmt.Errorf("حدث خطأ اثناء العثور على بيانات المستخدم")
	}

	if time.Now().After(userAuth.ExpiresAt) {
		return fmt.Errorf("مدة صلاحية رمز التحقق منتهية, يرجى طلب رمز جديد")
	}

	if userAuth.OTP != verifyOTPRequestDTO.OTP {
		userAuth.FailedAttempts++
		s.authRepo.UpdateAuth(userAuth)
		return fmt.Errorf("رمز التحقق المدخل غير صحيح")
	}

	userAuth.ExpiresAt = time.Now()
	userAuth.FailedAttempts = 0

	err = s.authRepo.UpdateAuth(userAuth)
	if err != nil {
		return fmt.Errorf("حدث خطأ أثناء حفظ بيانات التحقق")
	}
	return nil
}
