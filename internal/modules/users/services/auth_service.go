package users

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	authDTO "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/dtos"

	"github.com/Al-Khaimah/khaimah-golang-backend/internal/base"
	userModel "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
	userRepository "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
	"google.golang.org/api/idtoken"
)

type Provider interface {
	Exchange(idToken string) (string, error)
	GetClientSecret() (string, error)
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

type AuthService struct {
	userRepo  *userRepository.UserRepository
	authRepo  *userRepository.AuthRepository
	providers map[string]Provider
	jwtSecret []byte
}

func NewAuthService(userRepo *userRepository.UserRepository, authRepo *userRepository.AuthRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
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

func (g *GoogleProvider) Exchange(token string) (string, error) {
	p, err := idtoken.Validate(context.Background(), token, g.ClientID)
	if err != nil {
		return "", fmt.Errorf("google token invalid: %w", err)
	}
	email, _ := p.Claims["email"].(string)
	return email, nil
}

func (a *AppleProvider) Exchange(token string) (string, error) {
	set, err := jwk.Fetch(context.Background(), a.JWKSUrl)
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

func (g *GoogleProvider) GetClientSecret() (string, error) {
	return "", nil
}

func (a *AppleProvider) GetClientSecret() (string, error) {
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

func (s *AuthService) SSOLogin(dto *authDTO.OAuthRequestDTO, token string) base.Response {
	provider, ok := s.providers[dto.Provider]
	if !ok {
		return base.SetErrorMessage("unsupported provider")
	}

	email, err := provider.Exchange(token)
	if err != nil {
		return base.SetErrorMessage("invalid token: " + err.Error())
	}

	user, err := s.userRepo.FindOneByEmail(email)
	if err != nil {
		return base.SetErrorMessage("db error")
	}

	userExists := (user != nil)
	if !userExists {
		newUser := &userModel.User{
			Email: email,
		}
		user, err = s.userRepo.CreateUser(newUser)
		if err != nil {
			return base.SetErrorMessage("failed to create user")
		}
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(s.jwtSecret)
	if err != nil {
		return base.SetErrorMessage("sign token error")
	}

	var userDetails *authDTO.UserDetailsDTO
	if userExists {
		userDetails = &authDTO.UserDetailsDTO{
			ID:         user.ID.String(),
			FirstName:  user.FirstName,
			Email:      user.Email,
			Categories: user.Categories,
		}
	}

	SSOLoginResponse := authDTO.SSOLoginResponse{
		Token:      signed,
		UserExists: userExists,
		User:       userDetails,
	}

	return base.SetData(SSOLoginResponse, "login successful")
}

func (s *AuthService) ValidateWithProvider(providerKey, token string) (bool, error) { //ToDo :: Function is not used ?!
	provider, ok := s.providers[providerKey]
	if !ok {
		return false, fmt.Errorf("unsupported provider")
	}

	if providerKey == "apple" {
		clientSecret, err := provider.GetClientSecret()
		if err != nil {
			return false, fmt.Errorf("failed to generate client secret: %w", err)
		}
		_ = clientSecret
	}

	if providerKey == "google" {
		//ToDO :: Google Auth
	}

	_, err := provider.Exchange(token)
	return err == nil, err
}
