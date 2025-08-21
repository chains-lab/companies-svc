package handlers

import (
	"context"
	"fmt"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) BlockDistributor(ctx context.Context, req *disProto.BlockDistributorRequest) (*disProto.DistributorBlock, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	if initiator.Role != roles.Admin && initiator.Role != roles.SuperUser {
		s.Log(ctx).Warnf("user %s with role %s attempted to block distributor %s", initiator.ID, initiator.Role, req.DistributorId)

		return nil, errx.RaiseNoPermissions(
			ctx,
			fmt.Errorf("user %s with role %s does not have permission to block distributors", initiator.ID, initiator.Role),
		)
	}

	distributorID, err := uuid.Parse(req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID format: %s", req.DistributorId)

		return nil, errx.RaiseInvalidArgument(
			ctx,
			fmt.Errorf("invalid distributor ID format: %s", req.DistributorId),
			&errdetails.BadRequest_FieldViolation{
				Field:       "distributor_id",
				Description: "invalid UUID format for distributor ID",
			},
		)
	}

	block, err := s.app.BlockDistributor(ctx, initiator.ID, distributorID, req.Reason)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to block distributor %s", req.DistributorId)

		return nil, err
	}

	s.Log(ctx).Infof("distributor %s blocked successfully by user %s", req.DistributorId, initiator.ID)

	return responses.Block(block), nil
}
