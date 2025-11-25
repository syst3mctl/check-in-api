package service

import (
	"context"
	"errors"
	"time"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"
)

type AttendanceService struct {
	repo port.AttendanceRepository
}

func NewAttendanceService(repo port.AttendanceRepository) *AttendanceService {
	return &AttendanceService{repo: repo}
}

func (s *AttendanceService) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *AttendanceService) CheckIn(ctx context.Context, userID, orgID string, req *domain.Attendance) (*domain.Attendance, error) {
	// Check if already checked in
	latest, err := s.repo.GetLatestAttendance(ctx, userID)
	if err == nil && latest != nil && latest.CheckOutTime == nil {
		return nil, errors.New("already checked in")
	}

	req.UserID = userID
	req.OrgID = orgID
	req.CheckInTime = time.Now()
	req.Status = "PRESENT" // Default

	if req.TaskID != nil {
		// Task-based attendance
		task, err := s.repo.GetTaskByID(ctx, *req.TaskID)
		if err != nil {
			return nil, err
		}
		if task == nil {
			return nil, errors.New("task not found")
		}
		// TODO: Validate location (Geofencing)
		req.Type = "TASK"
	} else {
		// General attendance
		req.Type = "GENERAL"
		group, shift, err := s.repo.GetMemberGroup(ctx, orgID, userID)
		if err != nil {
			return nil, err
		}
		if group == nil {
			return nil, errors.New("user not in any group")
		}
		if shift != nil {
			req.ShiftApplied = shift.Name
			// Validate Shift Time (Simplified logic)
			// In a real app, parse shift.StartTime and compare with current time in shift.Timezone
		}
	}

	if err := s.repo.CreateAttendance(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *AttendanceService) CheckOut(ctx context.Context, userID string) error {
	latest, err := s.repo.GetLatestAttendance(ctx, userID)
	if err != nil {
		return err
	}
	if latest == nil || latest.CheckOutTime != nil {
		return errors.New("not checked in")
	}

	now := time.Now()
	latest.CheckOutTime = &now
	return s.repo.UpdateAttendance(ctx, latest)
}
