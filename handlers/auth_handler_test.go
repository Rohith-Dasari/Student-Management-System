package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sms/handlers"
	"sms/mocks"
	"sms/models"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		method         string
		body           any
		mockSetup      func(mock *mocks.MockAuthServiceI)
		expectedStatus int
		expectToken    bool
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body: map[string]string{
				"email":    "test@example.com",
				"password": "Password123!",
			},
			mockSetup: func(mock *mocks.MockAuthServiceI) {
				mock.EXPECT().ValidateLogin(gomock.Any(), "test@example.com", "Password123!").
					Return(models.User{UserID: "123", Email: "test@example.com", Role: "faculty"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			body:           nil,
			mockSetup:      func(mock *mocks.MockAuthServiceI) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid json",
			method:         http.MethodPost,
			body:           "{bad json",
			mockSetup:      func(mock *mocks.MockAuthServiceI) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid credentials",
			method: http.MethodPost,
			body: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			mockSetup: func(mock *mocks.MockAuthServiceI) {
				mock.EXPECT().ValidateLogin(gomock.Any(), "test@example.com", "wrongpassword").
					Return(models.User{}, errors.New("invalid email or password"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			req := httptest.NewRequest(tt.method, "/login", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			mockService := mocks.NewMockAuthServiceI(ctrl)
			tt.mockSetup(mockService)

			handler := handlers.NewAuthHandler(mockService)

			handler.Login(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectToken {
				var res handlers.LoginResponse
				if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if res.Token == "" {
					t.Error("expected token in response, got empty string")
				}
			}
		})
	}
}

func TestAuthHandler_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		method         string
		body           any
		mockSetup      func(mock *mocks.MockAuthServiceI)
		expectedStatus int
		expectToken    bool
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body: map[string]string{
				"name":     "John Doe",
				"email":    "john@example.com",
				"password": "Password123!",
			},
			mockSetup: func(mock *mocks.MockAuthServiceI) {
				mock.EXPECT().Signup(gomock.Any(), "John Doe", "john@example.com", "Password123!").
					Return(models.User{UserID: "123", Email: "john@example.com", Role: "faculty"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			body:           nil,
			mockSetup:      func(mock *mocks.MockAuthServiceI) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "invalid json",
			method:         http.MethodPost,
			body:           "{bad json",
			mockSetup:      func(mock *mocks.MockAuthServiceI) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "missing fields",
			method: http.MethodPost,
			body: map[string]string{
				"name":     "",
				"email":    "john@example.com",
				"password": "",
			},
			mockSetup:      func(mock *mocks.MockAuthServiceI) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "email already exists",
			method: http.MethodPost,
			body: map[string]string{
				"name":     "John Doe",
				"email":    "john@example.com",
				"password": "Password123!",
			},
			mockSetup: func(mock *mocks.MockAuthServiceI) {
				mock.EXPECT().Signup(gomock.Any(), "John Doe", "john@example.com", "Password123!").
					Return(models.User{}, errors.New("email already in use"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				b, _ := json.Marshal(v)
				bodyBytes = b
			}

			req := httptest.NewRequest(tt.method, "/signup", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			mockService := mocks.NewMockAuthServiceI(ctrl)
			tt.mockSetup(mockService)

			handler := handlers.NewAuthHandler(mockService)
			handler.Signup(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectToken {
				var res handlers.SignupResponse
				if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if res.Token == "" {
					t.Error("expected token in response, got empty string")
				}
			}
		})
	}
}
