package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) SetDistributorStatusActive(ctx context.Context, req *disProto.SetDistributorStatusActiveRequest) (*disProto.Distributor, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	distributorID, err := requests.DistributorID(ctx, req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID: %s", req.DistributorId)

		return nil, err
	}

	distributor, err := s.app.SetDistributorStatusActive(ctx, initiator.ID, distributorID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to set distributor %s status to active", req.DistributorId)

		return nil, err
	}

	s.Log(ctx).Infof("distributor %s status set to active successfully", distributor.ID)

	return responses.Distributor(distributor), nil
}
