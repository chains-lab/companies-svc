package rest

import (
	"net/http"

	"github.com/chains-lab/distributors-svc/internal/api/rest/handlers"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/logium"
	"github.com/go-chi/chi/v5"
)

type Rest struct {
	server   *http.Server
	router   *chi.Mux
	handlers handlers.Service

	log logium.Logger
	cfg config.Config
}

func NewRest(cfg config.Config, log logium.Logger, app *app.App) Rest {
	logger := log.WithField("module", "api")
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	hands := handlers.NewService(cfg, logger, app)

	return Rest{
		handlers: hands,
		router:   router,
		server:   server,
		log:      logger,
		cfg:      cfg,
	}
}
