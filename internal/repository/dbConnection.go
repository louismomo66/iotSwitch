package repository

import (
	"iot_switch/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	// conf := config.LoadConfig()
	// dbURL := "host=" + conf.DBHost + " user=" + conf.DBUser + " password=" + conf.DBPassword + " dbname=" + conf.DBName + " port=" + conf.DBPort + " sslmode=require"
	dbURL := "host=localhost port=5434 user=postgres password=XNEHk9iSGp9GItlxVuXYfmbEiTyugBuZ dbname=iotSwitch sslmode=disable"
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err // This line will not be reached because log.Fatalf exits the program
	}

	// Auto-migrate models
	err = conn.AutoMigrate(&models.User{},&models.Device{},&models.Relay{},&models.Schedule{})
	if err != nil {
		log.Printf("Failed to auto-migrate: %v", err)
		return conn, err
	}

	return conn, nil
}
