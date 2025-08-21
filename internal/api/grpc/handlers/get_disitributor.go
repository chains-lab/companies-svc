package handlers

import (
	"context"
	"fmt"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetDistributor(ctx context.Context, req *disProto.GetDistributorRequest) (*disProto.Distributor, error) {
	distributorID, err := uuid.Parse(req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID: %s", req.DistributorId)

		return nil, errx.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid distributor ID format must be uuid: %s", req.DistributorId),
			&errdetails.BadRequest_FieldViolation{
				Field:       "distributor_id",
				Description: fmt.Sprintf("invalid distributor ID format must be uuid: %s", req.DistributorId),
			},
		)
	}

	distributor, err := s.app.GetDistributor(ctx, distributorID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get distributor with ID: %s", distributorID)

		return nil, err
	}

	return responses.Distributor(distributor), nil
}
