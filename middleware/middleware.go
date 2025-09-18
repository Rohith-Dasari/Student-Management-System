package middleware

import (
	"context"
	"errors"
	"net/http"
	"sms/constants"
	"sms/services"
	"strings"
)

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := services.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), constants.ContextUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, constants.ContextUserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, constants.ContextUserRoleKey, claims.Role)

		next(w, r.WithContext(ctx))
	}
}

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(constants.ContextUserIDKey).(string)
	if !ok {
		return "", errors.New("userID not found in context")
	}
	return userID, nil
}

func GetUserEmail(ctx context.Context) (string, error) {
	email, ok := ctx.Value(constants.ContextUserEmailKey).(string)
	if !ok {
		return "", errors.New("user Email not found in context")
	}
	return email, nil
}

func GetUserRole(ctx context.Context) (constants.Role, error) {
	role, ok := ctx.Value(constants.ContextUserRoleKey).(constants.Role)
	if !ok {
		return "", errors.New("user Role not found in context")
	}
	return role, nil
}
