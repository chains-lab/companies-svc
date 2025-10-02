package handlers

import (
	"github.com/chains-lab/distributors-svc/internal"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/logium"
)

type Adapter struct {
	app *app.App
	log logium.Logger
	cfg internal.Config
}

func NewAdapter(cfg internal.Config, log logium.Logger, a *app.App) Adapter {
	return Adapter{
		app: a,
		cfg: cfg,
		log: log,
	}
}
