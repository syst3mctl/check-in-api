package postgres

import (
	"context"
	"errors"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"

	"github.com/jackc/pgx/v5"
)

type AttendanceRepository struct {
	db *DB
}

func NewAttendanceRepository(db *DB) port.AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) CreateTask(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (org_id, title, assigned_user_id, geofencing_enabled, location_name, latitude, longitude, radius_meters)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	executor := r.db.GetExecutor(ctx)
	return executor.QueryRow(ctx, query, task.OrgID, task.Title, task.AssignedUserID, task.GeofencingEnabled, task.LocationName, task.Latitude, task.Longitude, task.RadiusMeters).
		Scan(&task.ID, &task.CreatedAt)
}

func (r *AttendanceRepository) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `SELECT id, org_id, title, assigned_user_id, geofencing_enabled, location_name, latitude, longitude, radius_meters, created_at FROM tasks WHERE id = $1`
	task := &domain.Task{}
	executor := r.db.GetExecutor(ctx)
	err := executor.QueryRow(ctx, query, id).Scan(&task.ID, &task.OrgID, &task.Title, &task.AssignedUserID, &task.GeofencingEnabled, &task.LocationName, &task.Latitude, &task.Longitude, &task.RadiusMeters, &task.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *AttendanceRepository) CreateAttendance(ctx context.Context, attendance *domain.Attendance) error {
	query := `
		INSERT INTO attendance (user_id, org_id, task_id, check_in_time, status, type, shift_applied, location_lat, location_long, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`
	executor := r.db.GetExecutor(ctx)
	return executor.QueryRow(ctx, query, attendance.UserID, attendance.OrgID, attendance.TaskID, attendance.CheckInTime, attendance.Status, attendance.Type, attendance.ShiftApplied, attendance.LocationLat, attendance.LocationLong, attendance.Note).
		Scan(&attendance.ID, &attendance.CreatedAt)
}

func (r *AttendanceRepository) UpdateAttendance(ctx context.Context, attendance *domain.Attendance) error {
	query := `
		UPDATE attendance
		SET check_out_time = $2
		WHERE id = $1
	`
	executor := r.db.GetExecutor(ctx)
	_, err := executor.Exec(ctx, query, attendance.ID, attendance.CheckOutTime)
	return err
}

func (r *AttendanceRepository) GetLatestAttendance(ctx context.Context, userID string) (*domain.Attendance, error) {
	query := `
		SELECT id, user_id, org_id, task_id, check_in_time, check_out_time, status, type, shift_applied, location_lat, location_long, note, created_at
		FROM attendance
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	att := &domain.Attendance{}
	executor := r.db.GetExecutor(ctx)
	err := executor.QueryRow(ctx, query, userID).Scan(&att.ID, &att.UserID, &att.OrgID, &att.TaskID, &att.CheckInTime, &att.CheckOutTime, &att.Status, &att.Type, &att.ShiftApplied, &att.LocationLat, &att.LocationLong, &att.Note, &att.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return att, nil
}

func (r *AttendanceRepository) GetMemberGroup(ctx context.Context, orgID, userID string) (*domain.Group, *domain.Shift, error) {
	// Join organization_members -> groups -> shifts
	query := `
		SELECT g.id, g.name, g.shift_id, s.id, s.name, s.start_time, s.end_time, s.timezone, s.allowed_late_minutes, s.working_days
		FROM organization_members om
		JOIN groups g ON om.group_id = g.id
		LEFT JOIN shifts s ON g.shift_id = s.id
		WHERE om.org_id = $1 AND om.user_id = $2
	`
	group := &domain.Group{}
	shift := &domain.Shift{}

	executor := r.db.GetExecutor(ctx)
	var shiftID, shiftName, shiftStartTime, shiftEndTime, shiftTimezone *string
	var shiftLate *int
	var shiftDays []string

	err := executor.QueryRow(ctx, query, orgID, userID).Scan(
		&group.ID, &group.Name, &group.ShiftID,
		&shiftID, &shiftName, &shiftStartTime, &shiftEndTime, &shiftTimezone, &shiftLate, &shiftDays,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	if shiftID != nil {
		shift.ID = *shiftID
		shift.Name = *shiftName
		shift.StartTime = *shiftStartTime
		shift.EndTime = *shiftEndTime
		shift.Timezone = *shiftTimezone
		shift.AllowedLateMinutes = *shiftLate
		shift.WorkingDays = shiftDays
		return group, shift, nil
	}

	return group, nil, nil
}
