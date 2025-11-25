package handler

import (
	"encoding/json"
	"net/http"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/service"
	"github.com/syst3mctl/check-in-api/internal/pkg/response"
	"github.com/syst3mctl/check-in-api/internal/pkg/validator"

	"github.com/go-chi/chi/v5"
)

type AttendanceHandler struct {
	svc *service.AttendanceService
}

func NewAttendanceHandler(svc *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

// CreateTask godoc
// @Summary Create a task
// @Description Create a new task for the organization
// @Tags Task
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param request body domain.Task true "Task Request"
// @Success 201 {object} domain.Task
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/tasks [post]
func (h *AttendanceHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	var req domain.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.OrgID = orgID

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	task, err := h.svc.CreateTask(r.Context(), &req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, task)
}

// CheckIn godoc
// @Summary Check-in
// @Description Perform a check-in (General or Task-based)
// @Tags Attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.CheckInRequest true "Check-In Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Router /attendance/check-in [post]
func (h *AttendanceHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	var req domain.CheckInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	att := &domain.Attendance{
		TaskID:       req.TaskID,
		LocationLat:  req.Latitude,
		LocationLong: req.Longitude,
		Note:         req.Note,
	}

	result, err := h.svc.CheckIn(r.Context(), userID, req.OrganizationID, att)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":          "success",
		"check_in_time":   result.CheckInTime,
		"attendance_type": result.Type,
		"shift_applied":   result.ShiftApplied,
		"is_late":         result.Status == "LATE",
	})
}

// CheckOut godoc
// @Summary Check-out
// @Description Perform a check-out
// @Tags Attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} domain.ErrorResponse "bad request"
// @Router /attendance/check-out [post]
func (h *AttendanceHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	if err := h.svc.CheckOut(r.Context(), userID); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
