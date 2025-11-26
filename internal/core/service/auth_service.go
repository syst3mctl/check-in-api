package service

import (
	"context"
	"errors"
	"time"

	"github.com/syst3mctl/check-in-api/internal/config"
	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   port.UserRepository
	config *config.Config
}

func NewAuthService(repo port.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, config: cfg}
}

func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	// Check if user exists
	existing, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		PhoneNumber:  req.PhoneNumber,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenPair, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if claims["type"] != "refresh" {
		return nil, errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user id in token")
	}

	// Verify user still exists
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Generate new Access Token
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	newAccessTokenString, err := newAccessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Generate new Refresh Token
	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	newRefreshTokenString, err := newRefreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  newAccessTokenString,
		RefreshToken: newRefreshTokenString,
	}, nil
}
