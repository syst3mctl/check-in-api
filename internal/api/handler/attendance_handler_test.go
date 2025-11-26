package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/service"
)

// MockAttendanceRepository is a mock implementation of port.AttendanceRepository
type MockAttendanceRepository struct {
	mock.Mock
}

func (m *MockAttendanceRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockAttendanceRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockAttendanceRepository) CreateAttendance(ctx context.Context, attendance *domain.Attendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepository) UpdateAttendance(ctx context.Context, attendance *domain.Attendance) error {
	args := m.Called(ctx, attendance)
	return args.Error(0)
}

func (m *MockAttendanceRepository) GetLatestAttendance(ctx context.Context, userID string) (*domain.Attendance, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attendance), args.Error(1)
}

func (m *MockAttendanceRepository) GetMemberGroup(ctx context.Context, orgID, userID string) (*domain.Group, *domain.Shift, error) {
	args := m.Called(ctx, orgID, userID)
	group := args.Get(0)
	shift := args.Get(1)

	var g *domain.Group
	if group != nil {
		g = group.(*domain.Group)
	}

	var s *domain.Shift
	if shift != nil {
		s = shift.(*domain.Shift)
	}

	return g, s, args.Error(2)
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name           string
		input          domain.Task
		mockSetup      func(*MockAttendanceRepository)
		expectedStatus int
	}{
		{
			name: "Success",
			input: domain.Task{
				Title: "Test Task",
			},
			mockSetup: func(m *MockAttendanceRepository) {
				m.On("CreateTask", mock.Anything, mock.MatchedBy(func(task *domain.Task) bool {
					return task.Title == "Test Task"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid Input",
			input: domain.Task{
				Title: "", // Required
			},
			mockSetup: func(m *MockAttendanceRepository) {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAttendanceRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewAttendanceService(mockRepo)
			handler := NewAttendanceHandler(svc)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/organizations/org-1/tasks", bytes.NewBuffer(body))

			r := chi.NewRouter()
			r.Post("/organizations/{org_id}/tasks", handler.CreateTask)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCheckIn(t *testing.T) {
	validOrgID := "123e4567-e89b-12d3-a456-426614174000"
	validUserID := "user-123"

	tests := []struct {
		name           string
		input          domain.CheckInRequest
		mockSetup      func(*MockAttendanceRepository)
		expectedStatus int
	}{
		{
			name: "Success - General CheckIn",
			input: domain.CheckInRequest{
				OrganizationID: validOrgID,
				Latitude:       10.0,
				Longitude:      20.0,
			},
			mockSetup: func(m *MockAttendanceRepository) {
				m.On("GetLatestAttendance", mock.Anything, validUserID).Return(nil, nil)
				m.On("GetMemberGroup", mock.Anything, validOrgID, validUserID).Return(&domain.Group{ID: "group-1"}, &domain.Shift{Name: "Morning"}, nil)
				m.On("CreateAttendance", mock.Anything, mock.MatchedBy(func(a *domain.Attendance) bool {
					return a.Type == "GENERAL" && a.ShiftApplied == "Morning"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Already Checked In",
			input: domain.CheckInRequest{
				OrganizationID: validOrgID,
				Latitude:       10.0,
				Longitude:      20.0,
			},
			mockSetup: func(m *MockAttendanceRepository) {
				m.On("GetLatestAttendance", mock.Anything, validUserID).Return(&domain.Attendance{CheckOutTime: nil}, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAttendanceRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewAttendanceService(mockRepo)
			handler := NewAttendanceHandler(svc)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/attendance/check-in", bytes.NewBuffer(body))
			ctx := context.WithValue(req.Context(), "user_id", validUserID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.CheckIn(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCheckOut(t *testing.T) {
	validUserID := "user-123"

	tests := []struct {
		name           string
		mockSetup      func(*MockAttendanceRepository)
		expectedStatus int
	}{
		{
			name: "Success",
			mockSetup: func(m *MockAttendanceRepository) {
				m.On("GetLatestAttendance", mock.Anything, validUserID).Return(&domain.Attendance{CheckOutTime: nil}, nil)
				m.On("UpdateAttendance", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Not Checked In",
			mockSetup: func(m *MockAttendanceRepository) {
				m.On("GetLatestAttendance", mock.Anything, validUserID).Return(nil, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAttendanceRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewAttendanceService(mockRepo)
			handler := NewAttendanceHandler(svc)

			req, _ := http.NewRequest("POST", "/attendance/check-out", nil)
			ctx := context.WithValue(req.Context(), "user_id", validUserID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.CheckOut(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
