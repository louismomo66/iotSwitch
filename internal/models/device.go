package models

import "gorm.io/gorm"

//  type Device struct{
// 	gorm.Model
// 	ESP32ID string `json:"esp32_id" gorm:"uniqueIndex"`
// 	ID        uint    `gorm:"primaryKey"`
// 	NumRelays int     `json:"num_relays"`
// 	Relays    []Relay `json:"relays" gorm:"foreignKey:ESP32ID"`
//  }
type Device struct {
	gorm.Model                // Includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	ESP32ID    string         `json:"esp32_id" gorm:"uniqueIndex;primaryKey"` // Make ESP32ID both a unique index and the primary key
	NumRelays  int            `json:"num_relays"` // Number of relays, informational
	Relays     []Relay        `json:"relays" gorm:"foreignKey:DeviceESP32ID"` // Specify the foreign key
}

//  type Relay struct {
// 	ID      int   `gorm:"primaryKey"`
// 	ESP32ID uint   `json:"esp32_id"`
// 	Pin     int    `json:"pin"`
// 	State    string `json:"state"`
// }
type Relay struct {
	ID            int   `gorm:"primaryKey"`   // Primary key
	ESP32ID string `json:"device_esp32_id" gorm:"foreignKey:DeviceESP32ID"` // Foreign key that points to ESP32ID on Device
	Pin           int    `json:"pin"`          // The physical pin number on the device
	State         string `json:"state"`        // The current state of the relay, e.g., "on" or "off"
}

