package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) GetDistributor(ctx context.Context, req *disProto.GetDistributorRequest) (*disProto.Distributor, error) {
	distributorID, err := requests.DistributorID(ctx, req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID format: %s", req.DistributorId)

		return nil, err
	}

	distributor, err := s.app.GetDistributor(ctx, distributorID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get distributor with ID: %s", distributorID)

		return nil, err
	}

	return responses.Distributor(distributor), nil
}
