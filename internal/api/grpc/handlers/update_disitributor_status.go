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

func (s Service) UpdateDistributorIcon(ctx context.Context, req *disProto.UpdateDistributorIconRequest) (*disProto.Distributor, error) {
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

	res, err := s.app.UpdateDistributorIcon(ctx, initiator.ID, distributorID, req.IconUrl)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to update distributor icon for ID: %s", distributorID)

		return nil, err
	}

	s.Log(ctx).Infof("distributor %s icon updated successfully to %s", distributorID, req.IconUrl)

	return responses.Distributor(res), nil
}
