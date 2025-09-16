package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sms/middleware"
	"sms/models"
	"sms/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetUserID(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, "u1", userID)

		email, err := middleware.GetUserEmail(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", email)

		role, err := middleware.GetUserRole(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, models.Role("faculty"), role)

		w.WriteHeader(http.StatusOK)
	})

	validToken, err := services.GenerateJWT("u1", "user@example.com", "faculty")
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid format header",
			authHeader:     "Token abc",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			handler := middleware.JWTAuth(nextHandler)
			handler.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}
func TestGetUserHelpers_NotFound(t *testing.T) {
	ctx := context.Background()

	_, err := middleware.GetUserID(ctx)
	assert.Error(t, err)
	assert.Equal(t, "userID not found in context", err.Error())

	_, err = middleware.GetUserEmail(ctx)
	assert.Error(t, err)
	assert.Equal(t, "user Email not found in context", err.Error())

	_, err = middleware.GetUserRole(ctx)
	assert.Error(t, err)
	assert.Equal(t, "user Role not found in context", err.Error())
}
