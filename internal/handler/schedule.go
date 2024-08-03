package handler

import (
	"encoding/json"
	"fmt"
	"iot_switch/internal/models"
	"iot_switch/internal/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ScheduleHandler struct {
	DB *gorm.DB
}

func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var schedule models.Schedule
	if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request payload")
		return
	}
 // As StartTime is already a time.Time object, we directly convert it to UTC.
 	schedule.StartTime = schedule.StartTime.UTC()

 // Check if the provided StartTime is in the past
 	if schedule.StartTime.Before(time.Now().UTC()) {
		utils.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("start time cannot be in the past"), "Start time cannot be in the past")
		return
	}
	// Ensure the relay exists
	var relay models.Relay
	if err := h.DB.First(&relay, schedule.RelayID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err, "Relay not found")
		return
	}

	if err := h.DB.Create(&schedule).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to create schedule")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

func (h *ScheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
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
// Convert updated StartTime to UTC and check if it's in the past
updatedStartTime := schedule.StartTime.UTC()
if updatedStartTime.Before(time.Now().UTC()) {
	utils.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("start time cannot be in the past"), "Start time cannot be in the past")
	return
}

// Ensure the updated StartTime is properly assigned back if valid
schedule.StartTime = updatedStartTime
	// Ensure the relay exists
	var relay models.Relay
	if err := h.DB.First(&relay, schedule.RelayID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err, "Relay not found")
		return
	}

	if err := h.DB.Save(&schedule).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to update schedule")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}

func (h *ScheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
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
	if err := h.DB.Delete(&models.Schedule{}, scheduleID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to delete schedule")
		return
	}else{
		log.Printf("Schedule with ID: %d deleted successfully", scheduleID)
	}
	response := map[string]string{
        "message": "Schedule deleted successfully",
    }
	utils.WriteJSON(w, http.StatusOK, response)
}
func (h *ScheduleHandler) ActivateSchedule(w http.ResponseWriter, r *http.Request) {
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

	// Ensure the relay exists
	var relay models.Relay
	if err := h.DB.First(&relay, schedule.RelayID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err, "Relay not found")
		return
	}

	// Activate the selected schedule
	schedule.Active = true
	if err := h.DB.Save(&schedule).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to activate schedule")
		return
	}

	// Check if the schedule start time is in the future
	now := time.Now()
	if schedule.StartTime.After(now) {
		// Schedule the relay state change when the schedule starts
		time.AfterFunc(schedule.StartTime.Sub(now), func() {
			relay.State = "on"
			h.DB.Save(&relay)
			// Schedule the relay state reset after the schedule duration
			time.AfterFunc(time.Duration(schedule.Duration)*time.Second, func() {
				relay.State = "off"
				h.DB.Save(&relay)
			})
		})
	} else {
		// Schedule the relay state change immediately
		relay.State = "on"
		h.DB.Save(&relay)
		// Schedule the relay state reset after the schedule duration
		time.AfterFunc(time.Duration(schedule.Duration)*time.Second, func() {
			relay.State = "off"
			h.DB.Save(&relay)
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}


func (h *ScheduleHandler) DeactivateSchedule(w http.ResponseWriter, r *http.Request) {
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

	// Ensure the relay exists
	var relay models.Relay
	if err := h.DB.First(&relay, schedule.RelayID).Error; err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, err, "Relay not found")
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

func (h *ScheduleHandler) GetAllSchedules(w http.ResponseWriter, r *http.Request) {
	var schedules []models.Schedule

	if err := h.DB.Find(&schedules).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to fetch schedules")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}