package server

import (
	"encoding/json"
	"github.com/stringintech/security-101/server/auth"
	"net/http"
	"time"

	"github.com/stringintech/security-101/model"
	"github.com/stringintech/security-101/store"
)

type RegisterRequest struct {
	FullName string `json:"full_name"`
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
	auth  *auth.JwtService
}

func NewUserHandler(store *store.UserStore, jwt *auth.JwtService) *UserHandler {
	return &UserHandler{
		store: store,
		auth:  jwt,
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

	if _, exists := h.store.GetUserByUsername(req.Username); exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := model.User{
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := h.store.Create(user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	Encode(w, user)
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

	user, exists := h.store.GetUserByUsername(req.Username)
	if !exists {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := h.auth.ComparePasswords(user.GetPassword(), req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	Encode(w, LoginResponse{
		Token: token,
		User:  user.(model.User), // safe
	})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	Encode(w, user)
}

func Encode(w http.ResponseWriter, payload any) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}
