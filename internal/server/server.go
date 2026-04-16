package server

import (
	"fmt"
	"go-final-project/internal/api"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	defaultPort = 7540
	defaultWeb  = "web"
)

type Config struct {
	Port              int
	WebDir            string
	ReadHeaderTimeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		Port:              defaultPort,
		WebDir:            defaultWeb,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func NewServer(cfg Config) *http.Server {
	r := chi.NewRouter()

	api.Init(r)

	r.Handle("/*", http.FileServer(http.Dir(cfg.WebDir)))

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           r,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}
}
