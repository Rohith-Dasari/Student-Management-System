package services_test

import (
	"strings"
	"testing"

	"sms/services"
)

func TestGenerateAndValidateJWT_Success(t *testing.T) {
	userID := "123"
	email := "test@example.com"
	role := "admin"

	token, err := services.GenerateJWT(userID, email, role)
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}
	if token == "" {
		t.Fatal("expected token string, got empty")
	}

	claims, err := services.ValidateJWT(token)
	if err != nil {
		t.Fatalf("expected no error validating token, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected UserID %v, got %v", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected Email %v, got %v", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected Role %v, got %v", role, claims.Role)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidToken := "this.is.not.a.valid.token"

	claims, err := services.ValidateJWT(invalidToken)
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
	if claims != nil {
		t.Fatal("expected claims to be nil for invalid token")
	}
}

func TestValidateJWT_WrongSigningMethod(t *testing.T) {
	token := strings.Join([]string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		"eyJ1c2VyX2lkIjoiMTIzIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwicm9sZSI6ImFkbWluIn0",
		"invalidsignature",
	}, ".")

	claims, err := services.ValidateJWT(token)
	if err == nil {
		t.Fatal("expected error for wrong signing method, got nil")
	}
	if claims != nil {
		t.Fatal("expected claims to be nil for wrong signing method")
	}
}
