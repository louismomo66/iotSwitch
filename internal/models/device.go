package models

import "gorm.io/gorm"

 type Device struct{
	gorm.Model
	ESP32ID string `json:"esp32_id" gorm:"uniqueIndex"`
	ID        uint    `gorm:"primaryKey"`
	NumRelays int     `json:"num_relays"`
	Relays    []Relay `json:"relays" gorm:"foreignKey:ESP32ID"`
 }
 type Relay struct {
	ID      int   `gorm:"primaryKey"`
	ESP32ID uint   `json:"esp32_id"`
	Pin     int    `json:"pin"`
	State    string `json:"state"`
}

