package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/service"
)

// MockReportRepository is a mock implementation of port.ReportRepository
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) GetGroupPerformance(ctx context.Context, groupID string) (*domain.GroupPerformanceReport, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GroupPerformanceReport), args.Error(1)
}

func TestGetGroupPerformance(t *testing.T) {
	tests := []struct {
		name           string
		groupID        string
		mockSetup      func(*MockReportRepository)
		expectedStatus int
	}{
		{
			name:    "Success",
			groupID: "group-123",
			mockSetup: func(m *MockReportRepository) {
				m.On("GetGroupPerformance", mock.Anything, "group-123").Return(&domain.GroupPerformanceReport{GroupName: "Test Group"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Error",
			groupID: "group-error",
			mockSetup: func(m *MockReportRepository) {
				m.On("GetGroupPerformance", mock.Anything, "group-error").Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockReportRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewReportService(mockRepo)
			handler := NewReportHandler(svc)

			r := chi.NewRouter()
			r.Get("/organizations/{org_id}/reports/groups/{group_id}", handler.GetGroupPerformance)

			req, _ := http.NewRequest("GET", "/organizations/org-1/reports/groups/"+tt.groupID, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
