package handler

import (
	"encoding/json"
	"errors"

	"iot_switch/internal/repository"
	"iot_switch/internal/service"
	"iot_switch/internal/utils"
	"net/http"

	"gorm.io/gorm"
)

type AuthHandler struct {
	authService    service.AuthService
	userRepository repository.UserRepository
}

func NewAuthHandler(authService service.AuthService, userRepository repository.UserRepository) *AuthHandler {
	return &AuthHandler{authService, userRepository}
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request")
		return
	}

	// Check if email exists
	if _, err := h.userRepository.GetUserEmail(req.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.WriteJSONError(w, http.StatusNotFound, nil, "Email not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Error checking email")
		return
	}

	// Generate OTP
	token, err := h.authService.GenerateOTP(req.Email)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Error generating OTP")
		return
	}

	// Send OTP via email
	if err := utils.SendEmail(req.Email, "Your OTP Code", "Your OTP is: "+token); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to send OTP")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request")
		return
	}

	if err := h.authService.VerifyOTP(req.Email, req.OTP); err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, err, "Invalid or expired OTP")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "OTP verified"})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request")
		return
	}

	if err := h.authService.ResetPassword(req.Email, req.NewPassword); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to reset password")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Password reset successfully"})
}
