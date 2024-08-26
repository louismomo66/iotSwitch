package main

import (
	"iot_switch/internal/config"
	"iot_switch/internal/handler"
	midelware "iot_switch/internal/middleware"
	"iot_switch/internal/repository"
	"iot_switch/internal/routes"
	"iot_switch/internal/service"
	"iot_switch/internal/utils"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	err := godotenv.Load()
	if err != nil {
		logger.Log("Error loading .env file: %v", err)
	}
	cfg := config.LoadConfig()
	db, err := repository.ConnectDB()
	if err != nil {
		logger.Log("msg", "Failed to start server", "err", err)
			os.Exit(1)
	}
	userRepository := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	deviceHandler := handler.NewDeviceController(deviceRepo)
	handler.InitOAuth( )
	otpManager := utils.NewOTPManager()
	authService := service.NewAuthService(userRepository, otpManager)
	authHandler := handler.NewAuthHandler(authService, userRepository)
	scheduleHandler := &handler.ScheduleHandler{DB: db}

	scheduleChecker := &service.ScheduleChecker{DB: db}
	go scheduleChecker.StartScheduleChecker()

	
	r := mux.NewRouter()
	routes.SetupRoutes(r, authHandler, cfg.JWTSecret,scheduleHandler,deviceHandler)
	// Wrap the router with the LoggingMiddleware
	loggedRouter := midelware.LoggingMiddleware(logger)(r)

	if err := http.ListenAndServe(":9000",
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		)(loggedRouter)); err != nil {
			logger.Log("msg", "Failed to start server", "err", err)
			os.Exit(1)
	}
}
// Mukasa9090