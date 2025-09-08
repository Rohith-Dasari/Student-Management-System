package userrepository

import "sms/models"

type UserRepositoryI interface {
	AddUser(id string, name, email, password string) error
	GetUserByEmailID(email string) (*models.User, error)
}
