package postgres

import (
	"context"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"
)

type ReportRepository struct {
	db *DB
}

func NewReportRepository(db *DB) port.ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetGroupPerformance(ctx context.Context, groupID string) (*domain.GroupPerformanceReport, error) {
	query := `
		SELECT g.name, COALESCE(s.name, 'No Shift'),
		       (SELECT COUNT(*) FROM attendance a 
		        JOIN organization_members om ON a.user_id = om.user_id 
		        WHERE om.group_id = g.id AND a.status = 'LATE') as late_count
		FROM groups g
		LEFT JOIN shifts s ON g.shift_id = s.id
		WHERE g.id = $1
	`
	report := &domain.GroupPerformanceReport{}
	executor := r.db.GetExecutor(ctx)
	err := executor.QueryRow(ctx, query, groupID).Scan(&report.GroupName, &report.AssignedShift, &report.TotalLateCheckins)
	if err != nil {
		return nil, err
	}
	report.AttendanceRate = "98%" // Placeholder logic
	return report, nil
}
