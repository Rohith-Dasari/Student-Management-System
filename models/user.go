package models

type Role string

const (
	Admin   Role = "admin"
	Faculty Role = "faculty"
)

type User struct {
	UserID   string
	Name     string
	Email    string
	Password string
	Role     Role
}
