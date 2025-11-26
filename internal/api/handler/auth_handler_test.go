package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syst3mctl/check-in-api/internal/config"
	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/service"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of port.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name           string
		input          domain.RegisterRequest
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			input: domain.RegisterRequest{
				FullName:    "John Doe",
				Email:       "john@example.com",
				Password:    "password123",
				PhoneNumber: "1234567890",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "john@example.com").Return(nil, nil)
				m.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.Email == "john@example.com" && u.FullName == "John Doe"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "User Already Exists",
			input: domain.RegisterRequest{
				FullName:    "Jane Doe",
				Email:       "jane@example.com",
				Password:    "password123",
				PhoneNumber: "0987654321",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "jane@example.com").Return(&domain.User{Email: "jane@example.com"}, nil)
			},
			expectedStatus: http.StatusInternalServerError, // Handler returns 500 on service error
		},
		{
			name: "Validation Error - Missing Email",
			input: domain.RegisterRequest{
				FullName:    "No Email",
				Password:    "password123",
				PhoneNumber: "1234567890",
			},
			mockSetup: func(m *MockUserRepository) {
				// No mock calls expected
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			cfg := &config.Config{JWTSecret: "testsecret"}
			authService := service.NewAuthService(mockRepo, cfg)
			authHandler := NewAuthHandler(authService)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			authHandler.Register(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	validUser := &domain.User{
		ID:           "user-123",
		Email:        "john@example.com",
		PasswordHash: string(hashedPassword),
	}

	tests := []struct {
		name           string
		input          domain.LoginRequest
		mockSetup      func(*MockUserRepository)
		expectedStatus int
	}{
		{
			name: "Success",
			input: domain.LoginRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "john@example.com").Return(validUser, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "User Not Found",
			input: domain.LoginRequest{
				Email:    "unknown@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "unknown@example.com").Return(nil, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid Password",
			input: domain.LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "john@example.com").Return(validUser, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Database Error",
			input: domain.LoginRequest{
				Email:    "error@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByEmail", mock.Anything, "error@example.com").Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			cfg := &config.Config{JWTSecret: "testsecret"}
			authService := service.NewAuthService(mockRepo, cfg)
			authHandler := NewAuthHandler(authService)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			authHandler.Login(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
