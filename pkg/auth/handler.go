package auth

import (
	"encoding/json"
	"net/http"

	"github.com/cheebz/go-pub/pkg/config"
)

// Handler registers all /auth/* routes onto the provided mux.
type Handler struct {
	service Service
	conf    config.Configuration
}

func NewHandler(service Service, conf config.Configuration) *Handler {
	return &Handler{service: service, conf: conf}
}

// currentUser is the JSON shape returned to clients.
type currentUser struct {
	Username string `json:"username"`
}

// GET /auth/
// Returns the currently authenticated user, refreshing tokens if needed.
func (h *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	claims, err := h.service.Authenticate(w, r)
	if err != nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentUser{Username: claims.Username})
}

// POST /auth/login
// Body: application/x-www-form-urlencoded — username, password
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := h.service.Login(w, username, password)
	if err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentUser{Username: user.Username})
}

// POST /auth/register
// Body: application/x-www-form-urlencoded — username, password
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if !h.conf.Register {
		http.Error(w, `{"error":"registration disabled"}`, http.StatusForbidden)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, `{"error":"username and password required"}`, http.StatusBadRequest)
		return
	}
	user, err := h.service.Register(username, password)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(currentUser{Username: user.Username})
}

// GET /auth/logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = h.service.Logout(w, r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}

// GET /auth/logoutAll
func (h *Handler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	_ = h.service.LogoutAll(w, r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out all sessions"})
}

// POST /auth/password
// Body: application/x-www-form-urlencoded — new-password
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	newPassword := r.FormValue("new-password")
	if newPassword == "" {
		http.Error(w, `{"error":"new-password required"}`, http.StatusBadRequest)
		return
	}
	if err := h.service.ChangePassword(r, newPassword); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "password changed"})
}
