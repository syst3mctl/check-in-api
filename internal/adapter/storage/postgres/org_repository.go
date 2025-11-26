package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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
	err := executor.QueryRow(ctx, query, org.Name, org.Email, org.DefaultLocationLat, org.DefaultLocationLong).
		Scan(&org.ID, &org.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "organizations_email_key" {
				return &domain.DuplicateError{Field: "email"}
			}
			// If we had name unique constraint:
			// if pgErr.ConstraintName == "organizations_name_key" {
			// 	return &domain.DuplicateError{Field: "name"}
			// }
		}
		return err
	}
	return nil
}

func (r *OrgRepository) GetOrganizationByID(ctx context.Context, id string) (*domain.Organization, error) {
	query := `
		SELECT id, name, email, default_location_lat, default_location_long, created_at
		FROM organizations
		WHERE id = $1
	`
	var org domain.Organization
	executor := r.db.GetExecutor(ctx)
	err := executor.QueryRow(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.Email, &org.DefaultLocationLat, &org.DefaultLocationLong, &org.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *OrgRepository) UpdateOrganization(ctx context.Context, org *domain.Organization) error {
	query := `
		UPDATE organizations
		SET name = $2, email = $3, default_location_lat = $4, default_location_long = $5
		WHERE id = $1
	`
	executor := r.db.GetExecutor(ctx)
	_, err := executor.Exec(ctx, query, org.ID, org.Name, org.Email, org.DefaultLocationLat, org.DefaultLocationLong)
	return err
}

func (r *OrgRepository) DeleteOrganization(ctx context.Context, id string) error {
	query := `DELETE FROM organizations WHERE id = $1`
	executor := r.db.GetExecutor(ctx)
	_, err := executor.Exec(ctx, query, id)
	return err
}

func (r *OrgRepository) ListOrganizations(ctx context.Context, userID string) ([]*domain.Organization, error) {
	query := `
		SELECT o.id, o.name, o.email, o.default_location_lat, o.default_location_long, o.created_at
		FROM organizations o
		JOIN organization_members om ON o.id = om.org_id
		WHERE om.user_id = $1
	`
	executor := r.db.GetExecutor(ctx)
	rows, err := executor.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*domain.Organization
	for rows.Next() {
		var org domain.Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.Email, &org.DefaultLocationLat, &org.DefaultLocationLong, &org.CreatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, &org)
	}
	return orgs, nil
}

func (r *OrgRepository) GetMember(ctx context.Context, orgID, userID string) (*domain.OrganizationMember, error) {
	query := `
		SELECT id, org_id, user_id, role, group_id, created_at
		FROM organization_members
		WHERE org_id = $1 AND user_id = $2
	`
	var member domain.OrganizationMember
	executor := r.db.GetExecutor(ctx)
	err := executor.QueryRow(ctx, query, orgID, userID).Scan(
		&member.ID, &member.OrgID, &member.UserID, &member.Role, &member.GroupID, &member.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &member, nil
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
func (r *OrgRepository) GetOrganizationMembers(ctx context.Context, orgID string) ([]*domain.OrganizationMemberDetail, error) {
	query := `
		SELECT 
			om.id, om.org_id, om.user_id, om.role, om.group_id, om.created_at,
			u.id, u.full_name, u.email, u.phone_number, u.created_at
		FROM organization_members om
		JOIN users u ON om.user_id = u.id
		WHERE om.org_id = $1
	`
	executor := r.db.GetExecutor(ctx)
	rows, err := executor.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.OrganizationMemberDetail
	for rows.Next() {
		var m domain.OrganizationMemberDetail
		if err := rows.Scan(
			&m.ID, &m.OrgID, &m.UserID, &m.Role, &m.GroupID, &m.CreatedAt,
			&m.User.ID, &m.User.FullName, &m.User.Email, &m.User.PhoneNumber, &m.User.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, nil
}

func (r *OrgRepository) UpdateOrganizationMember(ctx context.Context, member *domain.OrganizationMember) error {
	query := `
		UPDATE organization_members
		SET role = $3, group_id = $4
		WHERE org_id = $1 AND user_id = $2
	`
	executor := r.db.GetExecutor(ctx)
	_, err := executor.Exec(ctx, query, member.OrgID, member.UserID, member.Role, member.GroupID)
	return err
}

func (r *OrgRepository) RemoveOrganizationMember(ctx context.Context, orgID, userID string) error {
	query := `DELETE FROM organization_members WHERE org_id = $1 AND user_id = $2`
	executor := r.db.GetExecutor(ctx)
	_, err := executor.Exec(ctx, query, orgID, userID)
	return err
}
