package service

import (
	"errors"
	models "iot_switch/iotSwitchApp/internal/models"
	"iot_switch/iotSwitchApp/internal/repository"
	"iot_switch/iotSwitchApp/internal/utils"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignUp(user *models.User) error
	Login(email, password string) (string, error)
	GenerateOTP(email string) (string, error)
	VerifyOTP(email, otp string) error
	ResetPassword(email, newPassword string) error
}

type authService struct {
	userRepository repository.UserRepository
	jwtSecret      string
	OTPManager     *utils.OTPManager
}

func NewAuthService(userRepository repository.UserRepository, jwtSecret string, otpManager *utils.OTPManager) AuthService {
	return &authService{userRepository, jwtSecret, otpManager}
}

func (s *authService) SignUp(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.userRepository.CreateUser(user)
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) GenerateOTP(email string) (string, error) {
	otp, err := s.OTPManager.GenerateOTP(email, 5*time.Minute)
	if err != nil {
		return "", err
	}
	return otp, nil
}

func (s *authService) VerifyOTP(email, otp string) error {
	if s.OTPManager.VeryfyOTP(email, otp) {
		return nil
	}
	return errors.New("invalid or expired OTP")
}

func (s *authService) ResetPassword(email, newPassword string) error {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.userRepository.UpdatePasswordByEmail(user.Email, user.Password)
}
