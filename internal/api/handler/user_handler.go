package handler

import (
	"encoding/json"
	"net/http"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/service"
	"github.com/syst3mctl/check-in-api/internal/pkg/response"
	"github.com/syst3mctl/check-in-api/internal/pkg/validator"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetMe godoc
// @Summary Get current user profile
// @Description Get profile of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} domain.ErrorResponse "unauthorized"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update profile of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Update Request"
// @Success 200 {object} domain.User
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 401 {object} domain.ErrorResponse "unauthorized"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /update-profile [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	user, err := h.svc.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// Logout godoc
// @Summary Logout user
// @Description Logout the user (client-side token removal, server-side placeholder)
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}
