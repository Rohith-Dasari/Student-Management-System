package services_test

import (
	"context"
	"testing"

	mockrepo "sms/mocks"
	"sms/models"
	"sms/services"

	"github.com/golang/mock/gomock"
)

func TestValidateLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockUserRepositoryI(ctrl)
	authSvc := services.NewAuthService(mockRepo)

	email := "test@example.com"
	password := "Password123!"

	hashedPassword, _ := authSvc.HashPassword(password)

	user := &models.User{
		UserID:   "123",
		Name:     "Test User",
		Email:    email,
		Password: hashedPassword,
		Role:     "faculty",
	}

	mockRepo.EXPECT().GetUserByEmailID(email).Return(user, nil)

	ctx := context.Background()
	result, err := authSvc.ValidateLogin(ctx, email, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Email != email {
		t.Errorf("expected email %v, got %v", email, result.Email)
	}
}

func TestValidateLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockUserRepositoryI(ctrl)
	authSvc := services.NewAuthService(mockRepo)

	email := "notfound@example.com"

	mockRepo.EXPECT().GetUserByEmailID(email).Return(nil, nil)

	ctx := context.Background()
	_, err := authSvc.ValidateLogin(ctx, email, "password")
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected 'user not found' error, got %v", err)
	}
}

func TestValidateLogin_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockUserRepositoryI(ctrl)
	authSvc := services.NewAuthService(mockRepo)

	email := "test@example.com"
	password := "Password123!"
	hashedPassword, _ := authSvc.HashPassword(password)

	user := &models.User{
		UserID:   "123",
		Name:     "Test User",
		Email:    email,
		Password: hashedPassword,
		Role:     "faculty",
	}

	mockRepo.EXPECT().GetUserByEmailID(email).Return(user, nil)

	ctx := context.Background()
	_, err := authSvc.ValidateLogin(ctx, email, "WrongPassword")
	if err == nil || err.Error() != "invalid email or password" {
		t.Fatalf("expected 'invalid email or password', got %v", err)
	}
}

func TestSignup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockUserRepositoryI(ctrl)
	authSvc := services.NewAuthService(mockRepo)

	email := "newuser@example.com"
	password := "StrongPass123!"
	name := "New User"

	mockRepo.EXPECT().GetUserByEmailID(email).Return(nil, nil)
	mockRepo.EXPECT().AddUser(gomock.Any(), name, email, gomock.Any()).Return(nil)

	ctx := context.Background()
	user, err := authSvc.Signup(ctx, name, email, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "" && user.Name != name {
		t.Errorf("expected user name %v, got %v", name, user.Name)
	}
}

func TestSignup_EmailAlreadyUsed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mockrepo.NewMockUserRepositoryI(ctrl)
	authSvc := services.NewAuthService(mockRepo)

	email := "existing@example.com"
	user := &models.User{Email: email}

	mockRepo.EXPECT().GetUserByEmailID(email).Return(user, nil)

	ctx := context.Background()
	_, err := authSvc.Signup(ctx, "Name", email, "StrongPass123!")
	if err == nil || err.Error() != "email already in use" {
		t.Fatalf("expected 'email already in use', got %v", err)
	}
}

func TestIsValidEmail(t *testing.T) {
	authSvc := services.NewAuthService(nil)

	valid := "a@b.com"
	invalid := "not-an-email"

	if !authSvc.IsValidEmail(valid) {
		t.Errorf("expected %v to be valid", valid)
	}
	if authSvc.IsValidEmail(invalid) {
		t.Errorf("expected %v to be invalid", invalid)
	}
}

func TestIsValidPassword(t *testing.T) {
	authSvc := services.NewAuthService(nil)

	valid := "StrongPass123!"
	invalid := "weakpass"

	if !authSvc.IsValidPassword(valid) {
		t.Errorf("expected valid password to pass")
	}
	if authSvc.IsValidPassword(invalid) {
		t.Errorf("expected invalid password to fail")
	}
}
