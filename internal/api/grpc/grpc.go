package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/chains-lab/distributors-svc/internal/api/grpc/handlers"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/interceptors"
	"github.com/chains-lab/distributors-svc/internal/app"
	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/distributors-svc/pkg/logger"
	"google.golang.org/grpc"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
)

func Run(ctx context.Context, cfg config.Config, log logger.Logger, app *app.App) error {
	requestId := interceptors.RequestID()
	userAuth := interceptors.UserJwtAuth(cfg.JWT.User.AccessToken.SecretKey)
	serviceAuth := interceptors.ServiceJwtAuth(cfg.JWT.Service.SecretKey)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestId,
			serviceAuth,
			userAuth,
		),
	)

	service := handlers.NewService(app, cfg, log)

	disProto.RegisterDistributorServiceServer(grpcServer, service)
	empProto.RegisterEmployeeServiceServer(grpcServer, service)

	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	log.Infof("gRPC server listening on %s", lis.Addr())

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down gRPC server …")
		grpcServer.GracefulStop()
		return nil
	case err := <-serveErrCh:
		return fmt.Errorf("gRPC Serve() exited: %w", err)
	}
}
