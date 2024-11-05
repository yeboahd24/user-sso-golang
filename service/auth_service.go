package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	// "github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v5"

	"github.com/yeboahd24/user-sso/config"
	"github.com/yeboahd24/user-sso/model"
	"github.com/yeboahd24/user-sso/repository"
	"github.com/yeboahd24/user-sso/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	oauthConfig *oauth2.Config
	jwtSecret   string
}

func NewAuthService(userRepo *repository.UserRepository, config *config.Config) *AuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     config.OAuth.Google.ClientID,
		ClientSecret: config.OAuth.Google.ClientSecret,
		RedirectURL:  config.OAuth.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		userRepo:    userRepo,
		oauthConfig: oauthConfig,
	}
}

func (s *AuthService) RegisterUser(email, password string) error {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    email,
		Password: hashedPassword,
	}

	return s.userRepo.CreateUser(user)
}

func (s *AuthService) HandleGoogleCallback(code string) (*model.User, error) {
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	googleUser, err := s.getGoogleUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.FindByEmail(googleUser.Email)
	if err == nil {
		// User exists, update SSO info
		err = s.userRepo.UpdateSSOInfo(existingUser.ID, "google", googleUser.Email)
		if err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	// User doesn't exist
	return nil, errors.New("no account found with this email. please register first")
}

func (s *AuthService) getGoogleUserInfo(accessToken string) (*model.GoogleOAuthResponse, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleUser model.GoogleOAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, err
	}

	return &googleUser, nil
}

func (s *AuthService) ValidateLogin(email, password string) (*model.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !util.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetGoogleAuthURL returns the Google OAuth URL
func (s *AuthService) GetGoogleAuthURL() string {
	return s.oauthConfig.AuthCodeURL("state")
}

// GenerateTokens generates JWT tokens for the user
func (s *AuthService) GenerateTokens(user *model.User) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the JWT token
func (s *AuthService) ValidateToken(tokenString string) (*util.Claims, error) {
	claims := &util.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
