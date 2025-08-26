package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) GetDistributorBlock(ctx context.Context, req *disProto.GetDistributorBlockRequest) (*disProto.DistributorBlock, error) {
	blockID, err := requests.BlockID(ctx, req.BlockId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid block ID format: %s", req.BlockId)

		return nil, err
	}

	block, err := s.app.GetBlock(ctx, blockID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get distributor block for ID: %s", blockID)

		return nil, err
	}

	return responses.Block(block), nil
}
