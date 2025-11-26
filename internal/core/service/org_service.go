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

func (s *OrgService) GetOrganization(ctx context.Context, id string) (*domain.Organization, error) {
	return s.repo.GetOrganizationByID(ctx, id)
}

func (s *OrgService) UpdateOrganization(ctx context.Context, userID string, req *domain.Organization) error {
	// Check if user is OWNER
	member, err := s.repo.GetMember(ctx, req.ID, userID)
	if err != nil {
		return err
	}
	if member.Role != "OWNER" {
		return errors.New("only owner can update organization")
	}

	// Fetch existing org to preserve fields if needed, or just update
	// For now assuming req contains all fields or we merge.
	// The repository update updates specific fields.

	return s.repo.UpdateOrganization(ctx, req)
}

func (s *OrgService) DeleteOrganization(ctx context.Context, userID, orgID string) error {
	// Check if user is OWNER
	member, err := s.repo.GetMember(ctx, orgID, userID)
	if err != nil {
		return err
	}
	if member.Role != "OWNER" {
		return errors.New("only owner can delete organization")
	}

	return s.repo.DeleteOrganization(ctx, orgID)
}

func (s *OrgService) ListOrganizations(ctx context.Context, userID string) ([]*domain.Organization, error) {
	return s.repo.ListOrganizations(ctx, userID)
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
func (s *OrgService) GetEmployees(ctx context.Context, requesterUserID, orgID string) ([]*domain.OrganizationMemberDetail, error) {
	requester, err := s.repo.GetMember(ctx, orgID, requesterUserID)
	if err != nil {
		return nil, err
	}
	if requester.Role != "OWNER" && requester.Role != "MANAGER" {
		return nil, errors.New("unauthorized")
	}
	return s.repo.GetOrganizationMembers(ctx, orgID)
}

func (s *OrgService) UpdateEmployee(ctx context.Context, requesterUserID, orgID, targetUserID string, role string, groupID *string) error {
	requester, err := s.repo.GetMember(ctx, orgID, requesterUserID)
	if err != nil {
		return err
	}
	if requester.Role != "OWNER" && requester.Role != "MANAGER" {
		return errors.New("unauthorized")
	}

	target, err := s.repo.GetMember(ctx, orgID, targetUserID)
	if err != nil {
		return err
	}

	if requester.Role == "MANAGER" {
		if target.Role == "OWNER" || target.Role == "MANAGER" {
			return errors.New("manager cannot update owner or other managers")
		}
		if role == "OWNER" || role == "MANAGER" {
			return errors.New("manager cannot promote to owner or manager")
		}
	}

	if target.Role == "OWNER" && role != "OWNER" {
		return errors.New("cannot demote owner")
	}

	target.Role = role
	target.GroupID = groupID
	return s.repo.UpdateOrganizationMember(ctx, target)
}

func (s *OrgService) RemoveEmployee(ctx context.Context, requesterUserID, orgID, targetUserID string) error {
	requester, err := s.repo.GetMember(ctx, orgID, requesterUserID)
	if err != nil {
		return err
	}
	if requester.Role != "OWNER" && requester.Role != "MANAGER" {
		return errors.New("unauthorized")
	}

	target, err := s.repo.GetMember(ctx, orgID, targetUserID)
	if err != nil {
		return err
	}

	if requester.Role == "MANAGER" {
		if target.Role == "OWNER" || target.Role == "MANAGER" {
			return errors.New("manager cannot remove owner or other managers")
		}
	}

	if target.Role == "OWNER" {
		return errors.New("cannot remove owner")
	}

	return s.repo.RemoveOrganizationMember(ctx, orgID, targetUserID)
}
