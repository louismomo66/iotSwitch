package handler

import (
	"encoding/json"
	"iot_switch/iotSwitchApp/internal/models"
	"iot_switch/iotSwitchApp/internal/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Schedulehundler struct{
DB *gorm.DB
}

func (h *Schedulehundler)CreateSchedule(w http.ResponseWriter,r *http.Request){
	var schedule models.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil{
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request payload")
		return
	}

	if err := h.DB.Create(&schedule).Error;err != nil{
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to create schedule")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

func (h *Schedulehundler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scheduleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid schedule ID")
		return
	}

	var schedule models.Schedule
	if err := h.DB.First(&schedule, scheduleID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err, "Schedule not found")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request payload")
		return
	}

	if err := h.DB.Save(&schedule).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to update schedule")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}

func (h *Schedulehundler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scheduleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid schedule ID")
		return
	}

	if err := h.DB.Delete(&models.Schedule{}, scheduleID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to delete schedule")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Schedulehundler) ActivateSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scheduleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid schedule ID")
		return
	}

	var schedule models.Schedule
	if err := h.DB.First(&schedule, scheduleID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err, "Schedule not found")
		return
	}

	// Activate the selected schedule
	schedule.Active = true
	if err := h.DB.Save(&schedule).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to activate schedule")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}


func (h *Schedulehundler) DeactivateSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scheduleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid schedule ID")
		return
	}

	var schedule models.Schedule
	if err := h.DB.First(&schedule, scheduleID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err, "Schedule not found")
		return
	}

	// Deactivate the selected schedule
	schedule.Active = false
	if err := h.DB.Save(&schedule).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to deactivate schedule")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}