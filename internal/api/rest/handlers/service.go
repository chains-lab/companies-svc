package handlers

import (
	"net/http"

	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/logium"
)

type Service struct {
	app *app.App
	log logium.Logger
	cfg config.Config
}

func NewService(cfg config.Config, log logium.Logger, a *app.App) Service {
	return Service{
		app: a,
		cfg: cfg,
		log: log,
	}
}

func (s Service) Log(r *http.Request) logium.Logger {
	return s.log
}
