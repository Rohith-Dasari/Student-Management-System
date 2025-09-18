package models

import "sms/constants"

type User struct {
	UserID   string
	Name     string
	Email    string
	Password string
	Role     constants.Role
}
