// package repository

// import (
// 	"iot_switch/internal/models"
// 	"log"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func ConnectDB() (*gorm.DB, error) {
// 	// conf := config.LoadConfig()
// 	// dbURL := "host=" + conf.DBHost + " user=" + conf.DBUser + " password=" + conf.DBPassword + " dbname=" + conf.DBName + " port=" + conf.DBPort + " sslmode=require"
// 	// dbURL := "host=localhost port=5434 user=postgres password=XNEHk9iSGp9GItlxVuXYfmbEiTyugBuZ dbname=iotSwitch sslmode=disable"
// 	// dbURL :=  os.Getenv("DNS")
// 	dbURL := "host=postgres port=5432 user=postgres password=XNEHk9iSGp9GItlxVuXYfmbEiTyugBuZ dbname=iotSwitch sslmode=disable"

// 	conn, err := gorm.Open(postgres.New(postgres.Config{
// 		DSN:                  dbURL,
// 		PreferSimpleProtocol: true,
// 	}), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 		return nil, err // This line will not be reached because log.Fatalf exits the program
// 	}

// 	// Auto-migrate models
// 	err = conn.AutoMigrate(&models.User{},&models.Device{},&models.Relay{},&models.Schedule{})
// 	if err != nil {
// 		log.Printf("Failed to auto-migrate: %v", err)
// 		return conn, err
// 	}

//		return conn, nil
//	}
package repository

import (
	"log"
	"time"

	"iot_switch/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "host=postgres port=5432 user=postgres password=XNEHk9iSGp9GItlxVuXYfmbEiTyugBuZ dbname=iotSwitch sslmode=disable"
	maxAttempts := 5
	var db *gorm.DB
	var err error

	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Successfully connected to the database.")
			break
		}
		log.Printf("Attempt %d/%d failed: %s. Retrying in 5 seconds...", i, maxAttempts, err)
		time.Sleep(5 * time.Second) // wait before retrying
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxAttempts, err)
		return nil, err
	}

	if err = db.AutoMigrate(&models.User{}, &models.Device{}, &models.Relay{}, &models.Schedule{}); err != nil {
		log.Printf("Failed to auto-migrate: %v", err)
		return db, err
	}

	return db, nil
}
