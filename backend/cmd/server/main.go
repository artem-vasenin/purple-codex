package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/uptime-monitor/backend/internal/app"
	"github.com/example/uptime-monitor/backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		log.Printf("server listening on %s", cfg.Addr)
		if err := a.Server.ListenAndServe(); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Close(shutdown); err != nil {
		log.Fatal(err)
	}
}
