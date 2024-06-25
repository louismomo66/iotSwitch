package models

import (
	"time"

	"gorm.io/gorm"
)

 type Schedule struct {
	gorm.Model
	RelayID   uint      `json:"relay_id"`
	StartTime time.Time `json:"start_time"`
	Duration  int       `json:"duration"` // duration in minutes
	Active    bool      `json:"active"`
}

