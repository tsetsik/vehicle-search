package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tsetsik/vehicle-search/internal/cache"
	"github.com/tsetsik/vehicle-search/internal/config"
	"github.com/tsetsik/vehicle-search/internal/core"
	"github.com/tsetsik/vehicle-search/internal/ratelimiter"

	"github.com/joho/godotenv"
)

type Service struct {
	cfg config.Config
}

func NewService() (*Service, error) {
	_, b, _, _ := runtime.Caller(0)
	if err := godotenv.Load(filepath.Dir(b) + "/../../.env"); err != nil {
		panic(err.Error())
	}

	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	cfg, err := config.LoadConfig(host, port)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Service{
		cfg: *cfg,
	}, nil
}

func (s *Service) Config() config.Config {
	return s.cfg
}

func (s *Service) Start(ctx context.Context) error {
	// Validate the config
	if err := s.cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	mux := mux.NewRouter()

	cache := cache.NewLRUCacheStore(100)
	itemStore := core.NewStore[core.Item](cache)

	httpResolver := NewHttpResolver(ctx, logger, s.cfg, itemStore)

	mux.HandleFunc("/search", httpResolver.Search).Methods(http.MethodGet)
	mux.HandleFunc("/listings", httpResolver.Listings).Methods(http.MethodPost)

	mux.Use(rateLimitMiddleware)

	logger.Info("Starting server", slog.String("host", s.cfg.Host), slog.Int("port", s.cfg.Port))
	return http.ListenAndServe(s.cfg.Host+":"+strconv.Itoa(s.cfg.Port), mux)
}

func (s *Service) Stop() error {
	// Implement any necessary cleanup logic here
	return nil
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	rateLimiter := ratelimiter.NewRateLimiter(5, 10*time.Second) // 5 requests per 10 seconds

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !rateLimiter.Allow(ip) {
			errStr := fmt.Sprintf("Too Many Requests for %s", ip)
			http.Error(w, errStr, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
