package midelware

import (
	"encoding/json"
	"fmt"
	"iot_switch/iotSwitchApp/internal/config"
	"iot_switch/iotSwitchApp/internal/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		conf := config.LoadConfig()
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            Err := utils.SetError(nil, "No Token Found")
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(Err)
            return
        }

        tokenString := strings.Split(authHeader, "Bearer ")[1]
        fmt.Println("Received Token:", tokenString) // Debugging

        mySigningKey := []byte(conf.JWTSecret)
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return mySigningKey, nil
        })

        if err != nil {
            Err := utils.SetError(err, "Your Token has expired")
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(Err)
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            role := claims["role"].(string)
            r.Header.Set("Role", role)
            handler.ServeHTTP(w, r)
            return
        }

        Err := utils.SetError(nil, "Not Authorized")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(Err)
    }
}
