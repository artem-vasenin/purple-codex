package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/example/uptime-monitor/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const cookieName = "refresh_token"

var errCredentials = errors.New("invalid credentials")

type Repository interface {
	CreateUser(context.Context, *store.User) error
	FindUser(context.Context, string) (store.User, error)
	CreateSession(context.Context, *store.RefreshSession) error
	FindSession(context.Context, string) (store.RefreshSession, error)
	RotateSession(context.Context, uuid.UUID, *store.RefreshSession) error
	RevokeSession(context.Context, uuid.UUID) error
}
type Service struct {
	repo                  Repository
	secret                string
	accessTTL, refreshTTL time.Duration
}

func NewService(repo Repository, secret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{repo: repo, secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func normalize(email string) string { return strings.ToLower(strings.TrimSpace(email)) }
func (s *Service) RegisterRoutes(r chi.Router, secure bool) {
	r.Post("/auth/register", s.register(secure))
	r.Post("/auth/login", s.login(secure))
	r.Post("/auth/refresh", s.refresh(secure))
	r.Post("/auth/logout", s.logout(secure))
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func readCredentials(r *http.Request) (credentials, error) {
	var c credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil || normalize(c.Email) == "" || c.Password == "" {
		return c, errors.New("invalid request")
	}
	c.Email = normalize(c.Email)
	return c, nil
}
func (s *Service) register(secure bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := readCredentials(r)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid request"})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": "internal error"})
			return
		}
		u := store.User{ID: uuid.New(), Email: c.Email, PasswordHash: string(hash)}
		if err = s.repo.CreateUser(r.Context(), &u); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				writeJSON(w, 409, map[string]string{"error": "email already registered"})
				return
			}
			writeJSON(w, 500, map[string]string{"error": "internal error"})
			return
		}
		s.issue(w, r, u, secure)
	}
}
func (s *Service) login(secure bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := readCredentials(r)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid request"})
			return
		}
		u, findErr := s.repo.FindUser(r.Context(), c.Email)
		if findErr != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(c.Password)) != nil {
			writeJSON(w, 401, map[string]string{"error": errCredentials.Error()})
			return
		}
		s.issue(w, r, u, secure)
	}
}
func (s *Service) issue(w http.ResponseWriter, r *http.Request, u store.User, secure bool) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": u.ID.String(), "iat": now.Unix(), "exp": now.Add(s.accessTTL).Unix(), "jti": uuid.New().String()})
	access, err := token.SignedString([]byte(s.secret))
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "internal error"})
		return
	}
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		writeJSON(w, 500, map[string]string{"error": "internal error"})
		return
	}
	refresh := base64.RawURLEncoding.EncodeToString(raw)
	hash := sha256.Sum256([]byte(refresh))
	session := &store.RefreshSession{ID: uuid.New(), UserID: u.ID, TokenHash: base64.RawURLEncoding.EncodeToString(hash[:]), ExpiresAt: now.Add(s.refreshTTL)}
	if err = s.repo.CreateSession(r.Context(), session); err != nil {
		writeJSON(w, 500, map[string]string{"error": "internal error"})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: refresh, Path: "/auth", HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode, MaxAge: int(s.refreshTTL.Seconds())})
	writeJSON(w, 200, map[string]any{"accessToken": access, "expiresIn": int(s.accessTTL.Seconds())})
}
func hashRefresh(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
func (s *Service) refresh(secure bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil {
			writeJSON(w, 401, map[string]string{"error": "invalid refresh token"})
			return
		}
		old, err := s.repo.FindSession(r.Context(), hashRefresh(c.Value))
		if err != nil || old.ExpiresAt.Before(time.Now()) || old.RevokedAt != nil {
			writeJSON(w, 401, map[string]string{"error": "invalid refresh token"})
			return
		}
		raw := make([]byte, 32)
		if _, err = rand.Read(raw); err != nil {
			writeJSON(w, 500, map[string]string{"error": "internal error"})
			return
		}
		nextRaw := base64.RawURLEncoding.EncodeToString(raw)
		next := &store.RefreshSession{ID: uuid.New(), UserID: old.UserID, TokenHash: hashRefresh(nextRaw), ExpiresAt: time.Now().Add(s.refreshTTL)}
		if err = s.repo.RotateSession(r.Context(), old.ID, next); err != nil {
			writeJSON(w, 401, map[string]string{"error": "invalid refresh token"})
			return
		}
		u := store.User{ID: old.UserID}
		s.issueAccessOnly(w, u, secure, nextRaw)
	}
}
func (s *Service) issueAccessOnly(w http.ResponseWriter, u store.User, secure bool, refresh string) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": u.ID.String(), "iat": now.Unix(), "exp": now.Add(s.accessTTL).Unix(), "jti": uuid.New().String()})
	access, err := token.SignedString([]byte(s.secret))
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "internal error"})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: refresh, Path: "/auth", HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode, MaxAge: int(s.refreshTTL.Seconds())})
	writeJSON(w, 200, map[string]any{"accessToken": access, "expiresIn": int(s.accessTTL.Seconds())})
}
func (s *Service) logout(secure bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(cookieName); err == nil {
			if session, findErr := s.repo.FindSession(r.Context(), hashRefresh(c.Value)); findErr == nil {
				_ = s.repo.RevokeSession(r.Context(), session.ID)
			}
		}
		http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/auth", HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode, MaxAge: -1})
		writeJSON(w, 200, map[string]string{"status": "ok"})
	}
}
