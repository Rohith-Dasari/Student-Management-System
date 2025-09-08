package services

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"sms/models"
	userrepository "sms/respository/userRepository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	ur userrepository.UserRepositoryI
}

func NewAuthService(ur userrepository.UserRepositoryI) *AuthService {
	return &AuthService{ur: ur}
}

func (a *AuthService) ValidateLogin(ctx context.Context, email, password string) (models.User, error) {
	user, err := a.ur.GetUserByEmailID(email)

	if user == nil && err == nil {
		return models.User{}, errors.New("user not found")
	}
	if err != nil {
		return models.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.User{}, errors.New("invalid email or password")
	}

	return *user, nil
}

func (a *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func (a *AuthService) IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (a *AuthService) IsValidPassword(password string) bool {
	var (
		minLength = len(password) >= 12
		hasNumber = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasUpper  = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower  = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasSymbol = regexp.MustCompile(`[!@#$%^&*()\-+]`).MatchString(password)
	)

	return minLength && hasNumber && hasUpper && hasLower && hasSymbol
}

func (a *AuthService) Signup(ctx context.Context, name, email, password string) (models.User, error) {
	if !a.IsValidEmail(email) {
		return models.User{}, errors.New("invalid email format")
	}

	if user, _ := a.ur.GetUserByEmailID(email); user != nil {
		return models.User{}, errors.New("email already in use")
	}

	if !a.IsValidPassword(password) {
		return models.User{}, errors.New("password must be at least 12 characters long, and include uppercase, lowercase, number, and symbol")
	}

	hashedPassword, err := a.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}

	uuid := uuid.New().String()
	err = a.ur.AddUser(uuid, name, email, hashedPassword)
	if err != nil {
		return models.User{}, err
	}

	return models.User{Name: name, UserID: uuid, Role: "faculty"}, nil
}
