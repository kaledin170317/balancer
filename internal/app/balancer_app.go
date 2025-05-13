package app

import (
	"balancer/internal/adapters/api/rest/controllers"
	"balancer/internal/adapters/api/rest/middleware"
	"balancer/internal/adapters/api/rest/proxy"
	"balancer/internal/adapters/db/postgreSQL"
	"balancer/internal/balancer"
	"balancer/internal/balancer/cheker"
	"balancer/internal/config"
	"balancer/internal/domain/models"
	"balancer/internal/domain/usecases"
	"balancer/internal/logger"
	"balancer/internal/ratelimit"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	logger.Init(slog.LevelInfo)
	log := logger.Base()

	log.Info("starting balancer application")

	cfg, err := config.LoadFinalConfig()
	if err != nil {
		log.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	log.Info("config loaded", "listen", cfg.ListenAddr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := connectToDB(cfg.DatabaseDSN)

	clientRepo := postgreSQL.NewClientRepository(db)
	clientUC, err := usecases.NewClientUseCase(ctx, clientRepo)
	if err != nil {
		log.Error("failed to init client usecase", "err", err)
		os.Exit(1)
	}
	log.Info("client usecase initialized")

	backends := makeBackends(cfg.Backends)
	log.Info("backends initialized", "count", len(backends))

	checker := cheker.NewChecker(backends, cfg.HealthCheck.Interval)
	checker.StartCheck(ctx)
	log.Info("health checker started")

	bal := balancer.GetBalancer(cfg.Algorithm, backends)
	log.Info("balancer selected", "type", cfg.Algorithm)

	r := chi.NewRouter()

	limiter := ratelimit.NewLimiterFromUseCase(clientUC)
	r.Use(middleware.Middleware(limiter))

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req = req.WithContext(logger.WithContext(req.Context(), logger.Base()))
			next.ServeHTTP(w, req)
		})
	})

	clientCtrl := controllers.NewClientController(clientUC)
	r.Mount("/clients", clientCtrl.Routes())
	r.Handle("/*", proxy.NewBalancer(bal))

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: r,
	}

	go func() {
		log.Info("server starting", "addr", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server listen error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
	} else {
		log.Info("server shutdown completed")
	}
}

func connectToDB(dsn string) *sqlx.DB {
	log := logger.Base()
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Error("database connection failed", "dsn", dsn, "err", err)
		os.Exit(1)
	}
	log.Info("connected to database")
	return db
}

func makeBackends(urls []string) []*models.Backend {
	var backends []*models.Backend
	for _, raw := range urls {
		u, err := url.Parse(raw)
		if err != nil {
			logger.Base().Error("invalid backend URL", "url", raw, "err", err)
			continue
		}
		backends = append(backends, models.NewBackend(u))
	}
	return backends
}
