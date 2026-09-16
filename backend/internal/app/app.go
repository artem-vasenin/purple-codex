package app

import (
	"context"
	"net/http"

	"github.com/example/uptime-monitor/backend/internal/auth"
	"github.com/example/uptime-monitor/backend/internal/config"
	"github.com/example/uptime-monitor/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type App struct {
	DB     *store.DB
	Server *http.Server
}

func New(cfg config.Config) (*App, error) {
	db, err := store.Open(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	service := auth.NewService(store.NewRepository(db.Gorm), cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	r := chi.NewRouter()
	service.RegisterRoutes(r, cfg.CookieSecure)
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		if err := db.SQL.PingContext(req.Context()); err != nil {
			http.Error(w, "unhealthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return &App{DB: db, Server: &http.Server{Addr: cfg.Addr, Handler: r}}, nil
}
func (a *App) Close(ctx context.Context) error { _ = a.Server.Shutdown(ctx); return a.DB.SQL.Close() }
