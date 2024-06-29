// package midelware

// import (
// 	"encoding/json"
// 	"fmt"
// 	"iot_switch/iotSwitchApp/internal/config"
// 	"iot_switch/iotSwitchApp/internal/utils"
// 	"net/http"
// 	"strings"

// 	"github.com/golang-jwt/jwt/v5"
// )

// func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
// 		conf := config.LoadConfig()
//         authHeader := r.Header.Get("Authorization")
//         if authHeader == "" {
//             Err := utils.SetError(nil, "No Token Found")
//             w.Header().Set("Content-Type", "application/json")
//             json.NewEncoder(w).Encode(Err)
//             return
//         }

//         tokenString := strings.Split(authHeader, "Bearer ")[1]
//         fmt.Println("Received Token:", tokenString) // Debugging

//         mySigningKey := []byte(conf.JWTSecret)
//         token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//             if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//                 return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//             }
//             return mySigningKey, nil
//         })

//         if err != nil {
//             Err := utils.SetError(err, "Your Token has expired")
//             w.Header().Set("Content-Type", "application/json")
//             json.NewEncoder(w).Encode(Err)
//             return
//         }

//         if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//             role := claims["role"].(string)
//             r.Header.Set("Role", role)
//             handler.ServeHTTP(w, r)
//             return
//         }

//	        Err := utils.SetError(nil, "Not Authorized")
//	        w.Header().Set("Content-Type", "application/json")
//	        json.NewEncoder(w).Encode(Err)
//	    }
//	}
package midelware

import (
	"encoding/json"
	"iot_switch/iotSwitchApp/internal/config"
	"iot_switch/iotSwitchApp/internal/utils"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/golang-jwt/jwt"
)

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
    cfg := config.LoadConfig()
    logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			
			Err := utils.SetError(nil, "No Token Found")
			json.NewEncoder(w).Encode(Err)
			return
		}

		var mySigningKey = []byte(cfg.JWTSecret)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, logger.Log("there was an error in parsing token.")
			}
			return mySigningKey, nil
		})

		if err != nil {
			Err := utils.SetError(err, "Your Token has been expired.")
			json.NewEncoder(w).Encode(Err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            role := claims["role"]
            if role == "admin" {
                r.Header.Set("Role", "admin")
            } else {
                r.Header.Set("Role", "user")
            }
            handler.ServeHTTP(w, r)
            return
        }
		
		reserr := utils.SetError(nil, "Not Authorized.")
		json.NewEncoder(w).Encode(reserr)
	}
}