package service

import (
    "errors"
    "iot_switch/internal/models"
    "iot_switch/internal/repository"
    "iot_switch/internal/utils"
    "time"

    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    SignUp(user *models.User) error
    GenerateOTP(email string) (string, error)
    VerifyOTP(email, otp string) error
    ResetPassword(email, newPassword string) error
    GenerateJWT(email, role string) (string, error)
}

type authService struct {
    userRepository repository.UserRepository
    OTPManager     *utils.OTPManager
}

func NewAuthService(userRepository repository.UserRepository, otpManager *utils.OTPManager) AuthService {
    return &authService{userRepository, otpManager}
}

func (s *authService) SignUp(user *models.User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
    return s.userRepository.CreateUser(user)
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

func (s *authService) GenerateJWT(email, role string) (string, error) {
    return utils.GenerateJWT(email, role)
}
