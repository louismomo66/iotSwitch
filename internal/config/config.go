package config

// internal/config/config.go

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBHost                 string
	DBPort                 string
	JWTSecret              string
	SMTPUser               string
	SMTPPass               string
	SMTPHost               string
	SMTPPort               string
	GOOGLE_CLIENT_ID       string
	GOOGLE_CLIENT_SECRET   string
	APPLE_CLIENT_ID        string
	APPLE_TEAM_ID          string
	FACEBOOK_CLIENT_ID     string
	FACEBOOK_CLIENT_SECRET string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBName:                 os.Getenv("DB_NAME"),
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		JWTSecret:              os.Getenv("JWT_SECRET"),
		SMTPUser:               os.Getenv("SMTP_USER"),
		SMTPPass:               os.Getenv("SMTP_PASS"),
		SMTPHost:               os.Getenv("SMTP_HOST"),
		SMTPPort:               os.Getenv("SMTP_PORT"),
		GOOGLE_CLIENT_ID:       os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		APPLE_CLIENT_ID:        os.Getenv("APPLE_CLIENT_ID"),
		APPLE_TEAM_ID:          os.Getenv("APPLE_TEAM_ID"),
		FACEBOOK_CLIENT_ID:     os.Getenv("FACEBOOK_CLIENT_ID"),
		FACEBOOK_CLIENT_SECRET: os.Getenv("FACEBOOK_CLIENT_SECRET"),
	}
}
