package rest

import (
	"net/http"
	"time"

	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/logium"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	server *http.Server
	router *chi.Mux

	log logium.Logger
	cfg config.Config
}

func NewRest(cfg config.Config, log logium.Logger) Service {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:              cfg.Server.Port,
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return Service{
		router: router,
		server: server,
		log:    logger,
		cfg:    cfg,
	}
}
