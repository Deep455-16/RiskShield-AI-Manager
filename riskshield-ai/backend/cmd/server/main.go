package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riskshield-ai/backend/internal/api"
	"github.com/riskshield-ai/backend/internal/audit"
	"github.com/riskshield-ai/backend/internal/auth"
	"github.com/riskshield-ai/backend/internal/jobs"
	"github.com/riskshield-ai/backend/internal/observability"
	"github.com/riskshield-ai/backend/internal/policy"
	"github.com/riskshield-ai/backend/internal/risk"
	"github.com/riskshield-ai/backend/internal/store"
	"github.com/riskshield-ai/backend/internal/datasets"
	"github.com/riskshield-ai/backend/internal/governance"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://riskshield:riskshield@localhost:5432/riskshield?sslmode=disable"
	}

	db, err := store.NewDB(context.Background(), dbURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := store.RunMigrations(dbURL); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "riskshield-dev-secret-change-in-production-min-32-chars-long"
	}

	authSvc := auth.NewService(db, jwtSecret)
	auditSvc := audit.NewService(db)
	riskSvc := risk.NewService(db, auditSvc)
	policySvc := policy.NewService(db, riskSvc, auditSvc)
	jobSvc := jobs.NewService(db)
	govSvc := governance.NewService(db, auditSvc)

	registry := datasets.NewRegistry("../data")
	replayEngine := datasets.NewReplayEngine(registry, db, riskSvc, policySvc)

	// Seed demo data if DEMO_MODE is set
	if os.Getenv("DEMO_MODE") == "true" {
		if err := seedDemoData(db, authSvc); err != nil {
			logger.Warn("demo seed failed", "error", err)
		} else {
			logger.Info("demo data seeded successfully")
		}
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(observability.StructuredLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"riskshield-api","version":"1.0.0"}`))
	})
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(context.Background()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"not_ready"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ready"}`))
	})

	api.RegisterRoutes(r, db, authSvc, auditSvc, riskSvc, policySvc, jobSvc, jwtSecret, registry, replayEngine, govSvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	}
	logger.Info("server stopped")
}

func seedDemoData(db *store.DB, authSvc *auth.Service) error {
	ctx := context.Background()
	// Create demo org
	orgID, err := db.CreateOrganization(ctx, "Demo Fintech Corp", "demo-fintech")
	if err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	// Create demo admin
	_, err = authSvc.Register(ctx, orgID, "admin@riskshield.demo", "DemoAdmin123!", "Demo Admin")
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}
