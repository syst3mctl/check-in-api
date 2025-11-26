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

// GetOrganization godoc
// @Summary Get organization details
// @Description Get details of a specific organization
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Success 200 {object} domain.Organization
// @Failure 404 {object} domain.ErrorResponse "organization not found"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id} [get]
func (h *OrgHandler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	org, err := h.svc.GetOrganization(r.Context(), orgID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, org)
}

// UpdateOrganization godoc
// @Summary Update organization details
// @Description Update details of an organization (Owner only)
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param request body domain.Organization true "Organization Request"
// @Success 200
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 403 {object} domain.ErrorResponse "forbidden"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id} [put]
func (h *OrgHandler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	userID := r.Context().Value("user_id").(string)
	var req domain.Organization
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = orgID

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	if err := h.svc.UpdateOrganization(r.Context(), userID, &req); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// DeleteOrganization godoc
// @Summary Delete organization
// @Description Delete an organization (Owner only)
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Success 200
// @Failure 403 {object} domain.ErrorResponse "forbidden"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id} [delete]
func (h *OrgHandler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	userID := r.Context().Value("user_id").(string)

	if err := h.svc.DeleteOrganization(r.Context(), userID, orgID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// ListOrganizations godoc
// @Summary List organizations
// @Description List organizations where the user is a member
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} domain.Organization
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations [get]
func (h *OrgHandler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	orgs, err := h.svc.ListOrganizations(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.WriteJSON(w, http.StatusOK, orgs)
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

// GetEmployees godoc
// @Summary Get employees
// @Description Get list of employees in the organization (Owner/Manager only)
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Success 200 {array} domain.OrganizationMemberDetail
// @Failure 403 {object} domain.ErrorResponse "forbidden"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/employees [get]
func (h *OrgHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	userID := r.Context().Value("user_id").(string)

	employees, err := h.svc.GetEmployees(r.Context(), userID, orgID)
	if err != nil {
		if err.Error() == "unauthorized" {
			response.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, employees)
}

// UpdateEmployee godoc
// @Summary Update employee
// @Description Update employee role or group (Owner/Manager only)
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param user_id path string true "User ID"
// @Param request body domain.UpdateEmployeeRequest true "Update Request"
// @Success 200
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 403 {object} domain.ErrorResponse "forbidden"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/employees/{user_id} [put]
func (h *OrgHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	targetUserID := chi.URLParam(r, "user_id")
	userID := r.Context().Value("user_id").(string)

	var req domain.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	if err := h.svc.UpdateEmployee(r.Context(), userID, orgID, targetUserID, req.Role, req.GroupID); err != nil {
		if err.Error() == "unauthorized" || err.Error() == "manager cannot update owner or other managers" || err.Error() == "manager cannot promote to owner or manager" || err.Error() == "cannot demote owner" {
			response.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// RemoveEmployee godoc
// @Summary Remove employee
// @Description Remove an employee from the organization (Owner/Manager only)
// @Tags Organization
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param org_id path string true "Organization ID"
// @Param user_id path string true "User ID"
// @Success 200
// @Failure 403 {object} domain.ErrorResponse "forbidden"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /organizations/{org_id}/employees/{user_id} [delete]
func (h *OrgHandler) RemoveEmployee(w http.ResponseWriter, r *http.Request) {
	orgID := chi.URLParam(r, "org_id")
	targetUserID := chi.URLParam(r, "user_id")
	userID := r.Context().Value("user_id").(string)

	if err := h.svc.RemoveEmployee(r.Context(), userID, orgID, targetUserID); err != nil {
		if err.Error() == "unauthorized" || err.Error() == "manager cannot remove owner or other managers" || err.Error() == "cannot remove owner" {
			response.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
