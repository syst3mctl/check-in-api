package postgres

import (
	"context"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"
)

type OrgRepository struct {
	db *DB
}

func NewOrgRepository(db *DB) port.OrgRepository {
	return &OrgRepository{db: db}
}

func (r *OrgRepository) CreateOrganization(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (name, email, default_location_lat, default_location_long)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	executor := r.db.GetExecutor(ctx)
	return executor.QueryRow(ctx, query, org.Name, org.Email, org.DefaultLocationLat, org.DefaultLocationLong).
		Scan(&org.ID, &org.CreatedAt)
}

func (r *OrgRepository) AddMember(ctx context.Context, member *domain.OrganizationMember) error {
	query := `
		INSERT INTO organization_members (org_id, user_id, role, group_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	executor := r.db.GetExecutor(ctx)
	return executor.QueryRow(ctx, query, member.OrgID, member.UserID, member.Role, member.GroupID).
		Scan(&member.ID, &member.CreatedAt)
}

func (r *OrgRepository) CreateShift(ctx context.Context, shift *domain.Shift) error {
	query := `
		INSERT INTO shifts (org_id, name, start_time, end_time, timezone, allowed_late_minutes, working_days)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	executor := r.db.GetExecutor(ctx)
	return executor.QueryRow(ctx, query, shift.OrgID, shift.Name, shift.StartTime, shift.EndTime, shift.Timezone, shift.AllowedLateMinutes, shift.WorkingDays).
		Scan(&shift.ID, &shift.CreatedAt)
}

func (r *OrgRepository) CreateGroup(ctx context.Context, group *domain.Group) error {
	query := `
		INSERT INTO groups (org_id, name, shift_id, manager_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	executor := r.db.GetExecutor(ctx)
	return executor.QueryRow(ctx, query, group.OrgID, group.Name, group.ShiftID, group.ManagerID).
		Scan(&group.ID, &group.CreatedAt)
}

func (r *OrgRepository) UpdateMemberGroup(ctx context.Context, orgID, userID, groupID string) error {
	query := `
		UPDATE organization_members
		SET group_id = $3
		WHERE org_id = $1 AND user_id = $2
	`
	executor := r.db.GetExecutor(ctx)
	_, err := executor.Exec(ctx, query, orgID, userID, groupID)
	return err
}
