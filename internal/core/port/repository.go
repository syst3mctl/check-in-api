package port

import (
	"context"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

type OrgRepository interface {
	CreateOrganization(ctx context.Context, org *domain.Organization) error
	AddMember(ctx context.Context, member *domain.OrganizationMember) error
	CreateShift(ctx context.Context, shift *domain.Shift) error
	CreateGroup(ctx context.Context, group *domain.Group) error
	UpdateMemberGroup(ctx context.Context, orgID, userID, groupID string) error
}

type AttendanceRepository interface {
	CreateTask(ctx context.Context, task *domain.Task) error
	GetTaskByID(ctx context.Context, id string) (*domain.Task, error)
	CreateAttendance(ctx context.Context, attendance *domain.Attendance) error
	UpdateAttendance(ctx context.Context, attendance *domain.Attendance) error
	GetLatestAttendance(ctx context.Context, userID string) (*domain.Attendance, error)
	GetMemberGroup(ctx context.Context, orgID, userID string) (*domain.Group, *domain.Shift, error)
}

type ReportRepository interface {
	GetGroupPerformance(ctx context.Context, groupID string) (*domain.GroupPerformanceReport, error)
}
