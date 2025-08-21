package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/logium"
)

type Service struct {
	app *app.App
	cfg config.Config
	log logium.Logger

	disProto.UnimplementedDistributorServiceServer
	empProto.UnimplementedEmployeeServiceServer
}

func NewService(a *app.App, cfg config.Config, log logium.Logger) *Service {
	return &Service{
		app: a,
		cfg: cfg,
		log: log,
	}
}

func (s Service) Log(ctx context.Context) logium.Logger {
	return s.log.WithField("request_id", meta.RequestID(ctx))
}
