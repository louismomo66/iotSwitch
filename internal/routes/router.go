package routes

import (
	"iot_switch/iotSwitchApp/internal/handler"
	midelware "iot_switch/iotSwitchApp/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, authHandler *handler.AuthHandler, jwtSecret string,scheduleHandler *handler.ScheduleHandler, esp32Handler *handler.DeviceController) {
	r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/forgot-password", authHandler.ForgotPassword).Methods("POST")
	r.HandleFunc("/verify-otp", authHandler.VerifyOTP).Methods("POST")
	r.HandleFunc("/reset-password", authHandler.ResetPassword).Methods("POST")

	r.HandleFunc("/auth/google/login", handler.HandleGoogleLogin).Methods("GET")
	r.HandleFunc("/auth/google/callback",  authHandler.HandleGoogleCallback).Methods("GET")
	r.HandleFunc("/auth/facebook/login", handler.HandleFacebookLogin).Methods("GET")
	r.HandleFunc("/auth/facebook/callback",  authHandler.HandleFacebookCallback).Methods("GET")
	// r.HandleFunc("/auth/apple/login", handler.HandleAppleLogin).Methods("GET")
	// r.HandleFunc("/auth/apple/callback", handler.HandleAppleCallback).Methods("GET")


	// r.HandleFunc("/schedules", midelware.IsAuthorized(scheduleHandler.CreateSchedule)).Methods("POST")
	r.HandleFunc("/schedules", scheduleHandler.CreateSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}", scheduleHandler.UpdateSchedule).Methods("PUT")
	r.HandleFunc("/schedules/{id}", scheduleHandler.DeleteSchedule).Methods("DELETE")
	r.HandleFunc("/schedules/{id}/activate", scheduleHandler.ActivateSchedule).Methods("POST")
	r.HandleFunc("/schedules/{id}/deactivate", scheduleHandler.DeactivateSchedule).Methods("POST")
	r.HandleFunc("/relay-states/{esp32_id}", scheduleHandler.GetRelayStates).Methods("GET")
	r.HandleFunc("/schedules", scheduleHandler.GetAllSchedules).Methods("GET")

	r.HandleFunc("/register-esp32", esp32Handler.RegisterDevice).Methods("POST")
	r.HandleFunc("/devices", esp32Handler.GetAllDevices).Methods("Get")
	r.HandleFunc("/devices/set-relay", esp32Handler.SetRelayState).Methods("POST")
	r.HandleFunc("/devices/{esp32_id}/relays", esp32Handler.GetRelaysByESP32ID).Methods("GET")
	r.HandleFunc("/devices/{esp32_id}", midelware.IsAuthorized(esp32Handler.DeleteDevice)).Methods("DELETE")


	
}
