package service

import (
	"context"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/port"
)

type ReportService struct {
	repo port.ReportRepository
}

func NewReportService(repo port.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetGroupPerformance(ctx context.Context, groupID string) (*domain.GroupPerformanceReport, error) {
	return s.repo.GetGroupPerformance(ctx, groupID)
}
