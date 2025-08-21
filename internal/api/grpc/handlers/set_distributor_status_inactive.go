package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) SetDistributorStatusInactive(ctx context.Context, req *disProto.SetDistributorStatusInactiveRequest) (*disProto.Distributor, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	distributorID, err := uuid.Parse(req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID: %s", req.DistributorId)

		return nil, errx.RaiseInvalidArgument(
			ctx, err,
			&errdetails.BadRequest_FieldViolation{
				Field:       "distributor_id",
				Description: "invalid UUID format for distributor ID",
			},
		)
	}

	distributor, err := s.app.SetDistributorStatusInactive(ctx, initiator.ID, distributorID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to set distributor %s status to inactive", req.DistributorId)

		return nil, err
	}

	s.Log(ctx).Infof("distributor %s status set to inactive successfully", distributor.ID)

	return responses.Distributor(distributor), nil
}
