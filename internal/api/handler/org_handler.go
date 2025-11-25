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

type OrgHandler struct {
	svc *service.OrgService
}

func NewOrgHandler(svc *service.OrgService) *OrgHandler {
	return &OrgHandler{svc: svc}
}

// CreateOrganization godoc
// @Summary Create a new organization
// @Description Create a new organization and assign the creator as OWNER
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.Organization true "Organization Request"
// @Success 201 {object} domain.Organization
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations [post]
func (h *OrgHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	var req domain.Organization
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	org, err := h.svc.CreateOrganization(r.Context(), userID, &req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, org)
}

// InviteEmployee godoc
// @Summary Invite an employee
// @Description Invite an employee to the organization
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param request body domain.InviteEmployeeRequest true "Invite Request"
// @Success 200
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/invitations [post]
func (h *OrgHandler) InviteEmployee(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	var req domain.InviteEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	if err := h.svc.InviteEmployee(r.Context(), orgID, req.Email, req.Role, req.GroupID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// CreateShift godoc
// @Summary Create a shift
// @Description Create a new shift for the organization
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param request body domain.Shift true "Shift Request"
// @Success 201 {object} domain.Shift
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/shifts [post]
func (h *OrgHandler) CreateShift(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	var req domain.Shift
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.OrgID = orgID

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	shift, err := h.svc.CreateShift(r.Context(), &req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, shift)
}

// CreateGroup godoc
// @Summary Create a group
// @Description Create a new group/team for the organization
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param request body domain.Group true "Group Request"
// @Success 201 {object} domain.Group
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/groups [post]
func (h *OrgHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	var req domain.Group
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.OrgID = orgID

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	group, err := h.svc.CreateGroup(r.Context(), &req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, group)
}

// AssignUserToGroup godoc
// @Summary Assign user to group
// @Description Assign a user to a specific group within the organization
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param user_id path string true "User ID"
// @Param request body domain.AssignUserToGroupRequest true "Assign Request"
// @Success 200
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/members/{user_id} [put]
func (h *OrgHandler) AssignUserToGroup(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	userID := chi.URLParam(r, "user_id")
	var req domain.AssignUserToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	if err := h.svc.AssignUserToGroup(r.Context(), orgID, userID, req.GroupID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
