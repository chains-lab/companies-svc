package handlers

import (
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/logium"
)

type Adapter struct {
	app *app.App
	log logium.Logger
	cfg config.Config
}

func NewAdapter(cfg config.Config, log logium.Logger, a *app.App) Adapter {
	return Adapter{
		app: a,
		cfg: cfg,
		log: log,
	}
}
