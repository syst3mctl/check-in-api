package handler

import (
	"encoding/json"
	"net/http"

	"github.com/syst3mctl/check-in-api/internal/core/domain"
	"github.com/syst3mctl/check-in-api/internal/core/service"
	"github.com/syst3mctl/check-in-api/internal/pkg/response"
	"github.com/syst3mctl/check-in-api/internal/pkg/validator"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Register Request"
// @Success 201 {object} domain.User
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 500 {object} domain.ErrorResponse "internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	user, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 401 {object} domain.ErrorResponse "unauthorized"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	tokenPair, err := h.svc.Login(r.Context(), &req)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, tokenPair)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} domain.TokenPair
// @Failure 400 {object} domain.ErrorResponse "invalid request body or validation errors"
// @Failure 401 {object} domain.ErrorResponse "unauthorized"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errResp := validator.ValidateStruct(&req); errResp != nil {
		response.WriteValidationError(w, errResp)
		return
	}

	tokenPair, err := h.svc.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, tokenPair)
}
