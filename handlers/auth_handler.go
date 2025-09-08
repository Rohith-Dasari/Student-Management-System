package handlers

import (
	"encoding/json"
	"net/http"
	"sms/services"
	"sms/utils"
)

type AuthHandler struct {
	as services.AuthServiceI
}

func NewAuthHandler(as services.AuthServiceI) *AuthHandler {
	return &AuthHandler{as: as}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.CustomError(w, http.StatusMethodNotAllowed, "invalid method")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.CustomError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.as.ValidateLogin(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.CustomError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := services.GenerateJWT(user.UserID, user.Email, string(user.Role))
	if err != nil {
		utils.CustomError(w, 409, "Failed to generate token")
		return
	}

	res := LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.CustomError(w, 409, "Failed to Encode Response")
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.CustomError(w, http.StatusMethodNotAllowed, "invalid method")
		return
	}

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.CustomError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.CustomError(w, http.StatusBadRequest, "name, email and password can't be empty ")
		return
	}

	user, err := h.as.Signup(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		utils.CustomError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := services.GenerateJWT(user.UserID, user.Email, string(user.Role))
	if err != nil {
		utils.CustomError(w, 409, "Failed to generate token")
		return
	}

	res := SignupResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		utils.CustomError(w, 409, "Failed to Encode Response")
	}
}
