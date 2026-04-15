package api

import (
	"encoding/json"
	"net/http"
	"wishlists_project/internal/application/service"
	"wishlists_project/internal/domain/requests"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req requests.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Password == "" || req.Email == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := handler.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Password == "" || req.Email == "" || req.Name == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := handler.authService.Register(req.Email, req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
