package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) CreateDistributor(ctx context.Context, req *disProto.CreateDistributorRequest) (*disProto.Distributor, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	distributor, err := s.app.CreateDistributor(ctx, initiator.ID, req.Name)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to create distributor")

		return nil, err
	}

	s.Log(ctx).Infof("distributor %s created successfully", distributor.ID)

	return responses.Distributor(distributor), nil
}
