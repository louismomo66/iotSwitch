package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"iot_switch/internal/models"
	"iot_switch/internal/repository"
	"iot_switch/internal/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repository.NewMockUserRepository(ctrl)
	mockAuthService := service.NewMockAuthService(ctrl)

	handler := AuthHandler{
		authService:    mockAuthService,
		userRepository: mockUserRepo,
	}

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:        "Successful Sign Up",
			requestBody: `{
				"first_name":"John",
				"second_name":"Doe",
				"email":"test@example.com",
				"phone_number":"1234567890",
				"username":"johndoe",
				"password":"password123",
				"role":"user"
			}`,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserEmail("test@example.com").Return("", nil)
				mockAuthService.EXPECT().SignUp(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "Email Already Exists",
			requestBody: `{
				"first_name":"John",
				"second_name":"Doe",
				"email":"test@example.com",
				"phone_number":"1234567890",
				"username":"johndoe",
				"password":"password123",
				"role":"user"
			}`,
			mockSetup: func() {
				mockUserRepo.EXPECT().GetUserEmail("test@example.com").Return("test@example.com", nil)
			},
			expectedStatus: http.StatusConflict,
		},
	
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockSetup()

			req := httptest.NewRequest("POST", "/signup", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			http.HandlerFunc(handler.SignUp).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repository.NewMockUserRepository(ctrl)
	mockAuthService := service.NewMockAuthService(ctrl)

	handler := AuthHandler{
		authService:    mockAuthService,
		userRepository: mockUserRepo,
	}

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:        "Successful Login",
			requestBody: `{
				"email":"test@example.com",
				"password":"password123"
			}`,
			mockSetup: func() {
				user := models.User{Email: "test@example.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZo5i.uE5vxB.l6b.b/8deK1boFJf39mjNP2u"} // bcrypt hash of "password123"
				mockUserRepo.EXPECT().GetUserByEmail("test@example.com").Return(user, nil)
				mockAuthService.EXPECT().GenerateJWT("test@example.com", "").Return("valid-token", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Invalid Credentials",
			requestBody: `{
				"email":"test@example.com",
				"password":"wrongpassword"
			}`,
			mockSetup: func() {
				user := models.User{Email: "test@example.com", Password: "$2a$10$N9qo8uLOickgx2ZMRZo5i.uE5vxB.l6b.b/8deK1boFJf39mjNP2u"}
				mockUserRepo.EXPECT().GetUserByEmail("test@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mockSetup()

			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			http.HandlerFunc(handler.Login).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
