package routes

import (
	"iot_switch/iotSwitchApp/internal/handler"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, authHandler *handler.AuthHandler, jwtSecret string) {
	r.HandleFunc("/signup", authHandler.SignUp).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/forgot-password", authHandler.ForgotPassword).Methods("POST")
	r.HandleFunc("/verify-otp", authHandler.VerifyOTP).Methods("POST")
	r.HandleFunc("/reset-password", authHandler.ResetPassword).Methods("POST")

	// r.Use(middleware.JWTMiddleware(jwtSecret))
	r.HandleFunc("/auth/google/login", handler.HandleGoogleLogin).Methods("GET")
	r.HandleFunc("/auth/google/callback",  authHandler.HandleGoogleCallback).Methods("GET")
	r.HandleFunc("/auth/facebook/login", handler.HandleFacebookLogin).Methods("GET")
	r.HandleFunc("/auth/facebook/callback",  authHandler.HandleFacebookCallback).Methods("GET")
	// r.HandleFunc("/auth/apple/login", handler.HandleAppleLogin).Methods("GET")
	// r.HandleFunc("/auth/apple/callback", handler.HandleAppleCallback).Methods("GET")

}
