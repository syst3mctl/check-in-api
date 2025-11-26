package handler

import (
	"bytes"
	"context"
	"encoding/json"
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

// MockOrgRepository is a mock implementation of port.OrgRepository
type MockOrgRepository struct {
	mock.Mock
}

func (m *MockOrgRepository) CreateOrganization(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrgRepository) GetOrganizationByID(ctx context.Context, id string) (*domain.Organization, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func (m *MockOrgRepository) UpdateOrganization(ctx context.Context, org *domain.Organization) error {
	args := m.Called(ctx, org)
	return args.Error(0)
}

func (m *MockOrgRepository) DeleteOrganization(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrgRepository) ListOrganizations(ctx context.Context, userID string) ([]*domain.Organization, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Organization), args.Error(1)
}

func (m *MockOrgRepository) GetMember(ctx context.Context, orgID, userID string) (*domain.OrganizationMember, error) {
	args := m.Called(ctx, orgID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrganizationMember), args.Error(1)
}

func (m *MockOrgRepository) GetOrganizationMembers(ctx context.Context, orgID string) ([]*domain.OrganizationMemberDetail, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.OrganizationMemberDetail), args.Error(1)
}

func (m *MockOrgRepository) UpdateOrganizationMember(ctx context.Context, member *domain.OrganizationMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockOrgRepository) RemoveOrganizationMember(ctx context.Context, orgID, userID string) error {
	args := m.Called(ctx, orgID, userID)
	return args.Error(0)
}

func (m *MockOrgRepository) AddMember(ctx context.Context, member *domain.OrganizationMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockOrgRepository) CreateShift(ctx context.Context, shift *domain.Shift) error {
	args := m.Called(ctx, shift)
	return args.Error(0)
}

func (m *MockOrgRepository) CreateGroup(ctx context.Context, group *domain.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockOrgRepository) UpdateMemberGroup(ctx context.Context, orgID, userID, groupID string) error {
	args := m.Called(ctx, orgID, userID, groupID)
	return args.Error(0)
}

// MockTransactionManager is a mock implementation of port.TransactionManager
type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// Simply execute the function for testing purposes
	return fn(ctx)
}

func TestCreateOrganization(t *testing.T) {
	tests := []struct {
		name           string
		input          domain.Organization
		mockSetup      func(*MockOrgRepository)
		expectedStatus int
	}{
		{
			name: "Success",
			input: domain.Organization{
				Name:                "Test Org",
				Email:               "test@example.com",
				DefaultLocationLat:  10.0,
				DefaultLocationLong: 20.0,
			},
			mockSetup: func(m *MockOrgRepository) {
				m.On("CreateOrganization", mock.Anything, mock.MatchedBy(func(o *domain.Organization) bool {
					return o.Name == "Test Org"
				})).Return(nil)
				m.On("AddMember", mock.Anything, mock.MatchedBy(func(mem *domain.OrganizationMember) bool {
					return mem.Role == "OWNER"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid Input",
			input: domain.Organization{
				Name: "", // Missing name
			},
			mockSetup: func(m *MockOrgRepository) {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockOrgRepository)
			mockUserRepo := new(MockUserRepository)
			mockTxMgr := new(MockTransactionManager)
			tt.mockSetup(mockRepo)

			svc := service.NewOrgService(mockRepo, mockUserRepo, mockTxMgr)
			handler := NewOrgHandler(svc)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/organizations", bytes.NewBuffer(body))
			// Inject user_id into context
			ctx := context.WithValue(req.Context(), "user_id", "user-123")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.CreateOrganization(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetOrganization(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		mockSetup      func(*MockOrgRepository)
		expectedStatus int
	}{
		{
			name:  "Success",
			orgID: "org-123",
			mockSetup: func(m *MockOrgRepository) {
				m.On("GetOrganizationByID", mock.Anything, "org-123").Return(&domain.Organization{ID: "org-123", Name: "Test Org"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "Not Found",
			orgID: "org-404",
			mockSetup: func(m *MockOrgRepository) {
				m.On("GetOrganizationByID", mock.Anything, "org-404").Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusInternalServerError, // Handler currently returns 500 on error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockOrgRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewOrgService(mockRepo, nil, nil)
			handler := NewOrgHandler(svc)

			r := chi.NewRouter()
			r.Get("/organizations/{org_id}", handler.GetOrganization)

			req, _ := http.NewRequest("GET", "/organizations/"+tt.orgID, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateOrganization(t *testing.T) {
	tests := []struct {
		name           string
		orgID          string
		userID         string
		input          domain.Organization
		mockSetup      func(*MockOrgRepository)
		expectedStatus int
	}{
		{
			name:   "Success - Owner",
			orgID:  "org-123",
			userID: "user-owner",
			input: domain.Organization{
				Name:                "Updated Name",
				Email:               "updated@example.com",
				DefaultLocationLat:  10.0,
				DefaultLocationLong: 20.0,
			},
			mockSetup: func(m *MockOrgRepository) {
				m.On("GetMember", mock.Anything, "org-123", "user-owner").Return(&domain.OrganizationMember{Role: "OWNER"}, nil)
				m.On("UpdateOrganization", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Forbidden - Not Owner",
			orgID:  "org-123",
			userID: "user-member",
			input: domain.Organization{
				Name:                "Updated Name",
				Email:               "updated@example.com",
				DefaultLocationLat:  10.0,
				DefaultLocationLong: 20.0,
			},
			mockSetup: func(m *MockOrgRepository) {
				m.On("GetMember", mock.Anything, "org-123", "user-member").Return(&domain.OrganizationMember{Role: "MEMBER"}, nil)
			},
			expectedStatus: http.StatusInternalServerError, // Service returns error, handler maps to 500 currently for generic errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockOrgRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewOrgService(mockRepo, nil, nil)
			handler := NewOrgHandler(svc)

			r := chi.NewRouter()
			r.Put("/organizations/{org_id}", handler.UpdateOrganization)

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("PUT", "/organizations/"+tt.orgID, bytes.NewBuffer(body))
			ctx := context.WithValue(req.Context(), "user_id", tt.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListOrganizations(t *testing.T) {
	mockRepo := new(MockOrgRepository)
	mockRepo.On("ListOrganizations", mock.Anything, "user-123").Return([]*domain.Organization{{ID: "org-1"}}, nil)

	svc := service.NewOrgService(mockRepo, nil, nil)
	handler := NewOrgHandler(svc)

	req, _ := http.NewRequest("GET", "/organizations", nil)
	ctx := context.WithValue(req.Context(), "user_id", "user-123")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ListOrganizations(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}
