package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/pkg/pagination"
)

func (s Service) SelectDistributorBlocks(ctx context.Context, req *disProto.SelectDistributorBlocksRequest) (*disProto.DistributorBlocksList, error) {
	filters := map[string]any{}
	if req.Filters.Status != nil {
		filters["status"] = req.Filters.Status
	}
	if req.Filters.InitiatorId != nil {
		filters["initiator_id"] = req.Filters.InitiatorId
	}
	if req.Filters.DistributorId != nil {
		filters["distributor_id"] = req.Filters.DistributorId
	}

	blocks, pag, err := s.app.SelectBlockages(ctx, filters, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to select distributor blocks")

		return nil, err
	}

	return responses.BlocksList(blocks, pag), err
}
