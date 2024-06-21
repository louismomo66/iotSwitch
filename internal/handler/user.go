package handler

import (
	"encoding/json"
	"iot_switch/iotSwitchApp/internal/models"
	"iot_switch/iotSwitchApp/internal/utils"
	"net/http"
)

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request")
		return
	}

	if err := h.authService.SignUp(&user); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Failed to create ")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request")
		return
	}

	token, err := h.authService.Login(request.Email, request.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, nil, "Incorect Email or password")
		return
	}

	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
