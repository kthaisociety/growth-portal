package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/kthaisociety/growth-portal/internal/config"
	"github.com/kthaisociety/growth-portal/internal/database"
	"github.com/kthaisociety/growth-portal/internal/handlers"
	"github.com/kthaisociety/growth-portal/internal/lumaservice"
	"github.com/kthaisociety/growth-portal/ui"
)

func main() {
	// Best effort: in a container the variables are set directly and there's no .env.
	_ = godotenv.Load()

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Connected but holding nothing yet: tables arrive with the first feature that stores data.
	db, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	router := handlers.NewRouter(handlers.Deps{
		Members: lumaservice.New(cfg.LumaServiceURL, cfg.LumaServiceSecret),
		PingDB:  func(ctx context.Context) error { return database.Ping(ctx, db) },
		UI:      ui.Dist(),
	})

	srv := &http.Server{
		Addr:              net.JoinHostPort(cfg.Host, cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("growth-portal listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
