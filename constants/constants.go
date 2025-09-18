package constants

type Role string

const (
	Admin   Role = "admin"
	Faculty Role = "faculty"
)

type contextKey string

const (
	ContextUserIDKey    contextKey = "userID"
	ContextUserEmailKey contextKey = "userEmail"
	ContextUserRoleKey  contextKey = "userRole"
)
