package auth

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/example/uptime-monitor/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type fakeRepository struct {
	users    map[string]store.User
	sessions map[string]store.RefreshSession
}

func newFake() *fakeRepository {
	return &fakeRepository{users: map[string]store.User{}, sessions: map[string]store.RefreshSession{}}
}
func (f *fakeRepository) CreateUser(_ context.Context, u *store.User) error {
	if _, ok := f.users[u.Email]; ok {
		return errCredentials
	}
	f.users[u.Email] = *u
	return nil
}
func (f *fakeRepository) FindUser(_ context.Context, email string) (store.User, error) {
	u, ok := f.users[email]
	if !ok {
		return store.User{}, errCredentials
	}
	return u, nil
}
func (f *fakeRepository) CreateSession(_ context.Context, s *store.RefreshSession) error {
	f.sessions[s.TokenHash] = *s
	return nil
}
func (f *fakeRepository) FindSession(_ context.Context, hash string) (store.RefreshSession, error) {
	s, ok := f.sessions[hash]
	if !ok {
		return store.RefreshSession{}, errCredentials
	}
	return s, nil
}
func (f *fakeRepository) RotateSession(_ context.Context, id uuid.UUID, next *store.RefreshSession) error {
	for hash, s := range f.sessions {
		if s.ID == id {
			now := time.Now()
			s.RevokedAt = &now
			f.sessions[hash] = s
			f.sessions[next.TokenHash] = *next
			return nil
		}
	}
	return errCredentials
}
func (f *fakeRepository) RevokeSession(_ context.Context, id uuid.UUID) error {
	for hash, s := range f.sessions {
		if s.ID == id {
			now := time.Now()
			s.RevokedAt = &now
			f.sessions[hash] = s
		}
	}
	return nil
}

func TestNormalizeEmail(t *testing.T) {
	if got := normalize("  User@Example.COM "); got != "user@example.com" {
		t.Fatalf("got %q", got)
	}
}
func TestRegisterAndLogin(t *testing.T) {
	repo := newFake()
	service := NewService(repo, "secret", time.Minute, time.Hour)
	router := chi.NewRouter()
	service.RegisterRoutes(router, false)
	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{"email":"User@example.com","password":"pass"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("register status %d", rec.Code)
	}
	if repo.users["user@example.com"].PasswordHash == "pass" {
		t.Fatal("stored plaintext password")
	}
	cookie := rec.Result().Cookies()[0]
	login := httptest.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"wrong"}`))
	login.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, login)
	if loginRec.Code != 401 {
		t.Fatalf("login status %d", loginRec.Code)
	}
	_ = cookie
}
