package auth

import (
	"context"
	"encoding/json"
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
func (f *fakeRepository) FindUserByID(_ context.Context, id uuid.UUID) (store.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return store.User{}, errCredentials
}
func (f *fakeRepository) UpdateUserProfile(_ context.Context, id uuid.UUID, firstName, lastName, phone string) error {
	for email, u := range f.users {
		if u.ID == id {
			u.FirstName, u.LastName, u.Phone = firstName, lastName, phone
			f.users[email] = u
			return nil
		}
	}
	return errCredentials
}
func (f *fakeRepository) UpdateUserAvatar(_ context.Context, id uuid.UUID, data []byte, contentType string) error {
	for email, u := range f.users {
		if u.ID == id {
			u.AvatarData, u.AvatarType = data, contentType
			f.users[email] = u
			return nil
		}
	}
	return errCredentials
}
func (f *fakeRepository) FindUserAvatar(_ context.Context, id uuid.UUID) ([]byte, string, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u.AvatarData, u.AvatarType, nil
		}
	}
	return nil, "", errCredentials
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

func TestProfileReadAndUpdate(t *testing.T) {
	repo := newFake()
	service := NewService(repo, "secret", time.Minute, time.Hour)
	router := chi.NewRouter()
	service.RegisterRoutes(router, false)

	register := httptest.NewRequest("POST", "/auth/register", strings.NewReader(`{"email":"profile@example.com","password":"pass"}`))
	register.Header.Set("Content-Type", "application/json")
	registerRec := httptest.NewRecorder()
	router.ServeHTTP(registerRec, register)
	var authResponse struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(registerRec.Body).Decode(&authResponse); err != nil || authResponse.AccessToken == "" {
		t.Fatal("register did not return an access token")
	}

	update := httptest.NewRequest("PATCH", "/auth/profile", strings.NewReader(`{"firstName":"Ada","lastName":"Lovelace","phone":"+123456789"}`))
	update.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)
	update.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, update)
	if updateRec.Code != 200 {
		t.Fatalf("profile update status %d", updateRec.Code)
	}

	profile := httptest.NewRequest("GET", "/auth/profile", nil)
	profile.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)
	profileRec := httptest.NewRecorder()
	router.ServeHTTP(profileRec, profile)
	var got profilePayload
	if err := json.NewDecoder(profileRec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.FirstName != "Ada" || got.LastName != "Lovelace" || got.Phone != "+123456789" {
		t.Fatalf("unexpected profile: %+v", got)
	}
}
