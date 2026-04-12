package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service is the interface used by middleware and the auth handler.
type Service interface {
	// Authenticate validates the JWT from the request, refreshing if expired.
	// Sets new cookies on refresh. Returns claims on success.
	Authenticate(w http.ResponseWriter, r *http.Request) (*Claims, error)
	// GetClaims validates the JWT only — no refresh, no side effects.
	GetClaims(r *http.Request) (*Claims, error)
	Login(w http.ResponseWriter, username, password string) (*User, error)
	Register(username, password string) (*User, error)
	Logout(w http.ResponseWriter, r *http.Request) error
	LogoutAll(w http.ResponseWriter, r *http.Request) error
	ChangePassword(r *http.Request, newPassword string) error
}

type service struct {
	repo      Repository
	jwt       *JWTHelper
	maxAgeSec int
}

func NewService(repo Repository, jwt *JWTHelper, refreshMaxAgeSecs int) Service {
	return &service{
		repo:      repo,
		jwt:       jwt,
		maxAgeSec: refreshMaxAgeSecs,
	}
}

func (s *service) GetClaims(r *http.Request) (*Claims, error) {
	return s.jwt.ParseJWT(r)
}

func (s *service) Authenticate(w http.ResponseWriter, r *http.Request) (*Claims, error) {
	claims, err := s.jwt.ParseJWT(r)
	if err == nil {
		return claims, nil
	}
	// JWT invalid or expired — try refresh
	return s.refresh(w, r)
}

func (s *service) refresh(w http.ResponseWriter, r *http.Request) (*Claims, error) {
	refreshClaims, err := s.jwt.ParseRefresh(r)
	if err != nil {
		return nil, err
	}
	if err := s.repo.ValidateRefresh(refreshClaims.UserID, refreshClaims.ID); err != nil {
		return nil, err
	}
	user, err := s.repo.GetUserByID(refreshClaims.UserID)
	if err != nil {
		return nil, err
	}
	signed, claims, err := s.jwt.CreateJWT(user.ID, user.Username, user.UUID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.InvalidateRefresh(refreshClaims.ID); err != nil {
		log.Println("auth: failed to invalidate refresh token:", err)
	}
	newRefresh, jti, err := s.jwt.CreateRefresh(user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveRefresh(user.ID, jti, s.maxAgeSec); err != nil {
		return nil, err
	}
	s.jwt.SetJWTCookie(w, signed)
	s.jwt.SetRefreshCookie(w, newRefresh)
	return claims, nil
}

func (s *service) Login(w http.ResponseWriter, username, password string) (*User, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	signed, _, err := s.jwt.CreateJWT(user.ID, user.Username, user.UUID)
	if err != nil {
		return nil, err
	}
	refreshSigned, jti, err := s.jwt.CreateRefresh(user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveRefresh(user.ID, jti, s.maxAgeSec); err != nil {
		return nil, err
	}
	s.jwt.SetJWTCookie(w, signed)
	s.jwt.SetRefreshCookie(w, refreshSigned)
	return &user, nil
}

func (s *service) Register(username, password string) (*User, error) {
	if _, err := s.repo.GetUserByUsername(username); err == nil {
		return nil, errors.New("username already taken")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	userUUID := uuid.New().String()
	user, err := s.repo.CreateUser(username, string(hash), userUUID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) Logout(w http.ResponseWriter, r *http.Request) error {
	refreshClaims, err := s.jwt.ParseRefresh(r)
	if err == nil {
		_ = s.repo.InvalidateRefresh(refreshClaims.ID)
	}
	s.jwt.ClearCookies(w)
	return nil
}

func (s *service) LogoutAll(w http.ResponseWriter, r *http.Request) error {
	refreshClaims, err := s.jwt.ParseRefresh(r)
	if err == nil {
		_ = s.repo.DeleteAllRefresh(refreshClaims.UserID)
	}
	s.jwt.ClearCookies(w)
	return nil
}

func (s *service) ChangePassword(r *http.Request, newPassword string) error {
	claims, err := s.jwt.ParseJWT(r)
	if err != nil {
		return errors.New("unauthorized")
	}
	user, err := s.repo.GetUserByUsername(claims.Username)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(user.ID, string(hash))
}

// PurgeWorker periodically deletes expired refresh tokens.
type PurgeWorker struct {
	repo     Repository
	interval time.Duration
}

func NewPurgeWorker(repo Repository) *PurgeWorker {
	return &PurgeWorker{repo: repo, interval: 24 * time.Hour}
}

func (p *PurgeWorker) Start() {
	for {
		time.Sleep(p.interval)
		if err := p.repo.DeleteExpiredRefresh(); err != nil {
			log.Println("auth: failed to purge expired refresh tokens:", err)
		}
	}
}
