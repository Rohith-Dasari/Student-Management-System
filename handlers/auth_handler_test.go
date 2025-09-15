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

	"go.uber.org/mock/gomock"
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
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			mockService := mocks.NewMockAuthServiceI(ctrl)
			tt.mockSetup(mockService)

			handler := handlers.NewAuthHandler(mockService)

			handler.Login(w, req)
			// res:=w.Result()
			// defer res.Body.Close()

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectToken {
				var resp map[string]any
				err := json.Unmarshal((w.Body.Bytes()), &resp)
				if err != nil {
					t.Errorf("unmarshalling response body failed")
				}

				if _, ok := resp["data"]; !ok {
					t.Errorf("expected to have token in data field but data field not found")
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
				var resp map[string]any
				err := json.Unmarshal((w.Body.Bytes()), &resp)
				if err != nil {
					t.Errorf("unmarshalling response body failed")
				}

				if _, ok := resp["data"]; !ok {
					t.Errorf("expected to have token in data field but data field not found")
				}
			}
		})
	}
}
