package service

import (
	"context"
	"errors"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"
)

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *domain.RegisterRequest) (*domain.User, error) {
	// We reuse RegisterRequest for update for now, or create a specific UpdateUserRequest
	// Assuming we update full name and phone number. Email and Password updates might require separate flows or verification.

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	// If we want to allow email update, we should check uniqueness again.
	// For now, let's stick to simple profile update (name, phone).

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
