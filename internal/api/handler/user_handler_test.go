package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/service"
)

func TestGetMe(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserRepository)
		expectedStatus int
	}{
		{
			name:   "Success",
			userID: "user-123",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, "user-123").Return(&domain.User{ID: "user-123", Email: "test@example.com"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "User Not Found",
			userID: "user-404",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, "user-404").Return(nil, nil)
			},
			expectedStatus: http.StatusInternalServerError, // Service returns error if nil, handler returns 500
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo)
			handler := NewUserHandler(svc)

			req, _ := http.NewRequest("GET", "/me", nil)
			ctx := context.WithValue(req.Context(), "user_id", tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.GetMe(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		input          domain.RegisterRequest
		mockSetup      func(*MockUserRepository)
		expectedStatus int
	}{
		{
			name:   "Success",
			userID: "user-123",
			input: domain.RegisterRequest{
				FullName:    "Updated Name",
				PhoneNumber: "9876543210",
				Email:       "test@example.com", // Required by validator
				Password:    "password123",      // Required by validator
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, "user-123").Return(&domain.User{ID: "user-123", FullName: "Old Name"}, nil)
				m.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
					return u.FullName == "Updated Name" && u.PhoneNumber == "9876543210"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Validation Error",
			userID: "user-123",
			input: domain.RegisterRequest{
				FullName: "", // Invalid
			},
			mockSetup: func(m *MockUserRepository) {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo)
			handler := NewUserHandler(svc)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("PUT", "/update-profile", bytes.NewBuffer(body))
			ctx := context.WithValue(req.Context(), "user_id", tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.UpdateProfile(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	handler := NewUserHandler(nil)
	req, _ := http.NewRequest("POST", "/logout", nil)
	rr := httptest.NewRecorder()

	handler.Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
