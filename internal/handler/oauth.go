package handler

import (
	"context"
	"encoding/json"
	"iot_switch/iotSwitchApp/internal/config"
	"iot_switch/iotSwitchApp/internal/models"
	"iot_switch/iotSwitchApp/internal/utils"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var (
	googleOAuthConfig   *oauth2.Config
	facebookOAuthConfig *oauth2.Config
	
)

func InitOAuth() {
	conf := config.LoadConfig()
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:9000/auth/google/callback",
		ClientID:     conf.GOOGLE_CLIENT_ID,
		ClientSecret: conf.GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	facebookOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:9000/auth/facebook/callback",
		ClientID:     conf.FACEBOOK_CLIENT_ID,
		ClientSecret: conf.FACEBOOK_CLIENT_SECRET,
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	}

	
}


func  HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	if googleOAuthConfig == nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, nil, "OAuth configuration is not initialized")
		return
	}

	url := googleOAuthConfig.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func(h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    token, err := googleOAuthConfig.Exchange(context.Background(), code)
    if err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to exchange token")
        return
    }

    response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    if err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to get user info")
        return
    }
    defer response.Body.Close()

    var userInfo map[string]interface{}
    if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to decode user info")
        return
    }

	log.Printf("User Info: %+v\n", userInfo)

    email, emailOk := userInfo["email"].(string)
    if !emailOk {
        utils.WriteJSONError(w, http.StatusInternalServerError, nil, "Invalid user info")
        return
    }

    // Assuming email as the username for now
    user := &models.User{
        Email:     email,
        Username:  email,
        Role:      "", // Default role
        // FirstName and SecondName are not available in the response
    }

    // Check if the user already exists in the database
	// var existingUser *models.User
    existingUser, err := h.userRepository.GetUserEmail(email)
    if err != nil && err != gorm.ErrRecordNotFound {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to check existing user")
        return
    }
    
    if existingUser == "" {
        if err := h.userRepository.CreateUser(user); err != nil {
            utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to save user info")
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(userInfo); err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to encode user info")
        return
    }
}

func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	url := facebookOAuthConfig.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    token, err := facebookOAuthConfig.Exchange(context.Background(), code)
    if err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to exchange token")
        return
    }

    response, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
    if err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to get user info")
        return
    }
    defer response.Body.Close()

    var userInfo map[string]interface{}
    if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to decode user info")
        return
    }

    email, emailOk := userInfo["email"].(string)
    name, nameOk := userInfo["name"].(string)

    if !emailOk || !nameOk {
        utils.WriteJSONError(w, http.StatusInternalServerError, nil, "Invalid user info")
        return
    }

	existingUser, err := h.userRepository.GetUserEmail(email)
    if err != nil && err != gorm.ErrRecordNotFound {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to check existing user")
        return
    }

	if existingUser == "" {
        // Create new user
        user := &models.User{
            Email:     email,
            Username:  name,
            Role:      "", // Default role
            // FirstName and SecondName are not available in the response
        }
        if err := h.userRepository.CreateUser(user); err != nil {
            utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to save user info")
            return
        }

	}
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(userInfo); err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to encode user info")
        return
    }
}

// func HandleAppleLogin(w http.ResponseWriter, r *http.Request) {
// 	url := appleOAuthConfig.AuthCodeURL("randomstate")
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// func HandleAppleCallback(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")
// 	token, err := appleOAuthConfig.Exchange(context.Background(), code)
// 	if err != nil {
// 		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "Access Token: %s\n", token.AccessToken)
// }
