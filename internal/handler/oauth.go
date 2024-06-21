package handler

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"iot_switch/iotSwitchApp/internal/config"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

var (
	googleOAuthConfig   *oauth2.Config
	facebookOAuthConfig *oauth2.Config
	appleOAuthConfig    *oauth2.Config
	privateKey          *rsa.PrivateKey
)

func InitOAuth() {
	conf := config.LoadConfig()
	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     conf.GOOGLE_CLIENT_ID,
		ClientSecret: conf.GOOGLE_CLIENT_SECRET,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	facebookOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/facebook/callback",
		ClientID:     conf.FACEBOOK_CLIENT_ID,
		ClientSecret: conf.FACEBOOK_CLIENT_SECRET,
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	}

	appleOAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/apple/callback",
		ClientID:     conf.APPLE_CLIENT_ID,
		ClientSecret: generateAppleClientSecret(),
		Scopes:       []string{"name", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://appleid.apple.com/auth/authorize",
			TokenURL: "https://appleid.apple.com/auth/token",
		},
	}
}
func generateAppleClientSecret() string {
	claims := jwt.MapClaims{
		"iss": os.Getenv("APPLE_TEAM_ID"),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": os.Getenv("APPLE_CLIENT_ID"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	clientSecret, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatalf("Failed to sign client secret: %v", err)
	}

	return clientSecret
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Info: %+v\n", userInfo)
}
func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	url := facebookOAuthConfig.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := facebookOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Info: %+v\n", userInfo)
}

func HandleAppleLogin(w http.ResponseWriter, r *http.Request) {
	url := appleOAuthConfig.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleAppleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := appleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Access Token: %s\n", token.AccessToken)
}
