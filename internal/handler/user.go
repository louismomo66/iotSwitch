package handler

import (
	"encoding/json"
	"errors"
	"iot_switch/iotSwitchApp/internal/models"
	"iot_switch/iotSwitchApp/internal/utils"
	"net/http"

	"gorm.io/gorm"
)

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request")
		return
	}
	email, err := h.userRepository.GetUserEmail(user.Email);
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.WriteJSONError(w, http.StatusInternalServerError,  err,"error checking email")
		return
	}
	
	if email != "" {
		utils.WriteJSONError(w, http.StatusConflict, nil,"email already in use")
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
	user, err := h.userRepository.GetUserByEmail(request.Email)
	if err != nil {
		utils.WriteJSONError(w,http.StatusBadRequest,err,"username or password is incorrect")
		return 
	}
	check := utils.CheckPasswordHash(request.Password, user.Password)

	if !check {
		utils.WriteJSONError(w,http.StatusBadRequest,err,"username or password is incorrect")
		return 
	}

	validtoken, err := utils.GenerateJWT(user.Email, user.Role)
	if err != nil {
		utils.WriteJSONError(w, http.StatusUnauthorized, nil, "Incorect Email or password")
		return
	}
token := models.Token{
Email: user.Email,
Role: user.Role,
TokenString: validtoken,
}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
