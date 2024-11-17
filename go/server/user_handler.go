package handler

import (
	"encoding/json"
	"net/http"

	"github.com/stringintech/security-101/auth"
	"github.com/stringintech/security-101/model"
	"github.com/stringintech/security-101/store"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

type UserHandler struct {
	store *store.UserStore
	auth  *auth.Service
}

func NewUserHandler(store *store.UserStore, auth *auth.Service) *UserHandler {
	return &UserHandler{
		store: store,
		auth:  auth,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if _, exists := h.store.FindByUsername(req.Username); exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := model.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
	}

	if err := h.store.Save(user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, exists := h.store.FindByUsername(req.Username)
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := h.auth.ComparePasswords(user.PasswordHash, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		User:  user,
	})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
