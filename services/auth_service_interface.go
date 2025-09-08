package services

import (
	"context"
	"sms/models"
)

//go:generate mockgen -destination=../mocks/auth_service_mock.go -package=mocks -source=auth_service_interface.go
type AuthServiceI interface {
	ValidateLogin(ctx context.Context, email, password string) (models.User, error)
	Signup(ctx context.Context, name, email, password string) (models.User, error)
}
