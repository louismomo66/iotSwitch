package utils

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type OTP struct {
	Email     string
	Token     string
	ExpiresAt time.Time
}

type OTPManager struct {
	otps map[string]OTP
	mu   sync.Mutex
}

func NewOTPManager() *OTPManager {
	return &OTPManager{
		otps: make(map[string]OTP),
	}
}

func (m *OTPManager) GenerateOTP(email string, expiration time.Duration) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	token := fmt.Sprintf("%04d", r.Intn(10000))

	expiresAt := time.Now().Add(expiration)
	m.otps[email] = OTP{
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return token, nil

}

func (m *OTPManager) VeryfyOTP(email, token string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	otp, ok := m.otps[email]
	if !ok {
		return false
	}

	if time.Now().After(otp.ExpiresAt) {
		return false
	}

	return otp.Token == token
}
