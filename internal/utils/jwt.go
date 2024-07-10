package utils

import (
	"iot_switch/internal/config"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
  }

  func GenerateJWT(email, role string) (string, error) {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	
	cfg := config.LoadConfig()
	var mySigningKey = []byte(cfg.JWTSecret)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] =  time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(mySigningKey)



	if err != nil {
		logger.Log("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}