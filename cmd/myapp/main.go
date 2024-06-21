package main

import (
	"iot_switch/iotSwitchApp/internal/config"
	"iot_switch/iotSwitchApp/internal/handler"
	"iot_switch/iotSwitchApp/internal/repository"
	"iot_switch/iotSwitchApp/internal/routes"
	"iot_switch/iotSwitchApp/internal/service"
	"iot_switch/iotSwitchApp/internal/utils"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	userRepository := repository.NewUserRepository(db)
	otpManager := utils.NewOTPManager()
	authService := service.NewAuthService(userRepository, cfg.JWTSecret, otpManager)
	authHandler := handler.NewAuthHandler(authService, userRepository)

	r := mux.NewRouter()
	routes.SetupRoutes(r, authHandler, cfg.JWTSecret)

	if err := http.ListenAndServe(":9000",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(r)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
