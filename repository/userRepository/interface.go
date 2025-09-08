package userrepository

import "sms/models"

//go:generate mockgen -destination=../../mocks/user_repo_mock.go -package=mocks -source=interface.go
type UserRepositoryI interface {
	AddUser(id string, name, email, password string) error
	GetUserByEmailID(email string) (*models.User, error)
}
