package service

import (
	"context"
	"errors"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"
)

type OrgService struct {
	repo     port.OrgRepository
	userRepo port.UserRepository
	txMgr    port.TransactionManager
}

func NewOrgService(repo port.OrgRepository, userRepo port.UserRepository, txMgr port.TransactionManager) *OrgService {
	return &OrgService{repo: repo, userRepo: userRepo, txMgr: txMgr}
}

func (s *OrgService) CreateOrganization(ctx context.Context, userID string, req *domain.Organization) (*domain.Organization, error) {
	// Transactional: Create Org and Add User as Owner
	err := s.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.repo.CreateOrganization(ctx, req); err != nil {
			return err
		}

		member := &domain.OrganizationMember{
			OrgID:  req.ID,
			UserID: userID,
			Role:   "OWNER",
		}
		if err := s.repo.AddMember(ctx, member); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return req, nil
}

func (s *OrgService) CreateShift(ctx context.Context, shift *domain.Shift) (*domain.Shift, error) {
	if err := s.repo.CreateShift(ctx, shift); err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *OrgService) CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	if err := s.repo.CreateGroup(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *OrgService) InviteEmployee(ctx context.Context, orgID string, email string, role string, groupID *string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	member := &domain.OrganizationMember{
		OrgID:   orgID,
		UserID:  user.ID,
		Role:    role,
		GroupID: groupID,
	}
	return s.repo.AddMember(ctx, member)
}

func (s *OrgService) AssignUserToGroup(ctx context.Context, orgID, userID, groupID string) error {
	return s.repo.UpdateMemberGroup(ctx, orgID, userID, groupID)
}
