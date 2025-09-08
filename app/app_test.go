package app_test

import (
	"net/http"
	"net/http/httptest"
	"sms/app"
	"testing"
)

func TestSetupServerRoutes(t *testing.T) {
	db, err := app.InitDBWithDSN(":memory:")
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	mux := app.SetupServer(db)

	tests := []struct {
		method string
		route  string
	}{
		{"POST", "/api/v1/login"},
		{"POST", "/api/v1/signup"},
		{"POST", "/api/v1/students"},
		{"PATCH", "/api/v1/students/{studentID}"},
		{"POST", "/api/v1/grades"},
		{"GET", "/api/v1/grades"},
		{"GET", "/api/v1/grades/toppers"},
		{"PATCH", "/api/v1/grades"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.route, nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized && w.Code != http.StatusBadRequest {
			t.Errorf("route %s %s returned status %d", tt.method, tt.route, w.Code)
		}
	}
}
