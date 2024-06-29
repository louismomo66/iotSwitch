package repository

import (
	"iot_switch/iotSwitchApp/internal/models"
	"time"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetDevice(id string) (models.Device, error)
	CreateDevice(device *models.Device) error
	UpdateDevice(device *models.Device) error
	GetAllDevices() ([]models.Device, error)
	DeleteDeviceByESP32ID(esp32ID string) error
	UpdateRelayState(relay *models.Relay) error
	GetRelayState(relayID uint) string 
	GetDeviceByESP32ID(esp32ID string) (*models.Device, error)
	GetRelayByESP32IDAndPin(esp32ID string, pin int) (*models.Relay, error)
	GetRelayByESP32ID(esp32ID string) (*models.Device, error)
}

type GormDeviceRepo struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *GormDeviceRepo {
	return &GormDeviceRepo{db: db}
}

func (repo *GormDeviceRepo) GetDevice(id string) (models.Device, error) {
	var device models.Device
	err := repo.db.Preload("Relays").Where("esp32_id = ?", id).First(&device).Error
	return device, err
}

func (repo *GormDeviceRepo) CreateDevice(device *models.Device) error {
	return repo.db.Create(&device).Error
}

func (repo *GormDeviceRepo) UpdateDevice(device *models.Device) error {
	return repo.db.Save(&device).Error
}

func (repo *GormDeviceRepo) GetAllDevices() ([]models.Device, error) {
	var devices []models.Device
	err := repo.db.Find(&devices).Error
	return devices, err
}

func (repo *GormDeviceRepo) DeleteDeviceByESP32ID(esp32ID string) error {
    return repo.db.Where("esp32_id = ?", esp32ID).Delete(&models.Device{}).Error
}
func (repo *GormDeviceRepo) UpdateRelayState(relay *models.Relay) error {
	return repo.db.Model(&models.Relay{}).Where("id = ?", relay.ID).Update("state", relay.State).Error
}
func (repo *GormDeviceRepo) GetRelayState(relayID uint) string {
	var relay models.Relay
	if err := repo.db.First(&relay, relayID).Error; err != nil {
		return "unknown"
	}

	// Check if there's any active schedule
	var activeSchedule models.Schedule
	if err := repo.db.Where("relay_id = ? AND start_time <= ? AND end_time >= ?", relayID, time.Now(), time.Now()).First(&activeSchedule).Error; err == nil {
		return "on"
	}

	return relay.State
}
func (repo *GormDeviceRepo) GetDeviceByESP32ID(esp32ID string) (*models.Device, error) {
	var device models.Device
	
	err := repo.db.Preload("Relays").Where("esp32_id = ?", esp32ID).First(&device).Error
	return &device, err
}

func (repo *GormDeviceRepo) GetRelayByESP32IDAndPin(esp32ID string, pin int) (*models.Relay, error) {
	var relay models.Relay
	err := repo.db.Where("esp32_id = ? AND pin = ?", esp32ID, pin).First(&relay).Error
	return &relay, err
}
func (repo *GormDeviceRepo) GetRelayByESP32ID(esp32ID string) (*models.Device, error) {
    var device models.Device
    err := repo.db.Preload("Relays").Where("esp32_id = ?", esp32ID).First(&device).Error
    return &device, err
}
