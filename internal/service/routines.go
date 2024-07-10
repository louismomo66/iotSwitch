package service

import (
	"iot_switch/internal/models"
	"time"

	"gorm.io/gorm"
)
type ScheduleChecker struct {
    DB *gorm.DB
}

func (h *ScheduleChecker) StartScheduleChecker() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.checkSchedules()
	}
}

func (h *ScheduleChecker) checkSchedules() {
	now := time.Now()
	var schedules []models.Schedule

	// Fetch all active schedules whose duration has expired
	if err := h.DB.Where("active = ? AND start_time + interval '1 second' * duration < ?", true, now).Find(&schedules).Error; err != nil {
		return // Handle error as needed
	}

	for _, schedule := range schedules {
		schedule.Active = false
		h.DB.Save(&schedule) // Deactivate the schedule

		// Update the relay state to "off"
		var relay models.Relay
		if err := h.DB.First(&relay, schedule.RelayID).Error; err == nil {
			relay.State = "off"
			h.DB.Save(&relay)
		}
	}
}
