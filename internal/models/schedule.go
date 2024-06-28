package models

import (
	"time"

	"gorm.io/gorm"
)

 type Schedule struct {
	gorm.Model
	ESP32ID   string    `json:"esp32_id"`
	RelayID   int      `json:"relay_id"`
	StartTime time.Time `json:"start_time"`
	Duration  int       `json:"duration"` // duration in minutes
	Active    bool      `json:"active"`
}

