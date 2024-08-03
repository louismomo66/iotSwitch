package handler

import (
	"encoding/json"
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
 // Parse the StartTime to ensure it's valid and to adjust for timezone
 userTime, err := time.Parse(time.RFC3339, schedule.StartTime.Format(time.RFC3339))
 if err != nil {
	 utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid start time format")
	 return
 }

 // Convert user provided time to UTC (if server operates in UTC)
 schedule.StartTime = userTime.UTC()
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