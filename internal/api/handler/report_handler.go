package handler

import (
	"net/http"

	"github.com/syst3mctl/check-in-api/internal/core/service"
	"github.com/syst3mctl/check-in-api/internal/pkg/response"

	"github.com/go-chi/chi/v5"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// GetGroupPerformance godoc
// @Summary Get group performance report
// @Description Get performance report for a specific group
// @Tags Report
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param group_id path string true "Group ID"
// @Success 200 {object} domain.GroupPerformanceReport
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/reports/groups/{group_id} [get]
func (h *ReportHandler) GetGroupPerformance(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "group_id")
	report, err := h.svc.GetGroupPerformance(r.Context(), groupID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, report)
}
