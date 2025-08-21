package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetDistributorBlock(ctx context.Context, req *disProto.GetDistributorBlockRequest) (*disProto.DistributorBlock, error) {
	blockID, err := uuid.Parse(req.BlockId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid block ID: %s", req.BlockId)

		return nil, errx.RaiseInvalidArgument(
			ctx, err,
			&errdetails.BadRequest_FieldViolation{
				Field:       "block_id",
				Description: "invalid block ID format must be uuid",
			},
		)
	}

	block, err := s.app.GetBlock(ctx, blockID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get distributor block for ID: %s", blockID)

		return nil, err
	}

	return responses.Block(block), nil
}
