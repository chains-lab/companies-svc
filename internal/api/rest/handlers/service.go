package handlers

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/logium"
)

type Service struct {
	app app.App
	log logium.Logger
	cfg config.Config
}

func NewService(cfg config.Config, log logium.Logger, a app.App) Service {
	return Service{
		app: a,
		cfg: cfg,
		log: log,
	}
}

func (s Service) Log(ctx context.Context) logium.Logger {
	return s.log.WithField("request_id", meta.RequestID(ctx))
}
