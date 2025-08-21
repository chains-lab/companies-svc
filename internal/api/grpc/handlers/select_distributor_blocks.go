package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/pagi"
)

func (s Service) SelectDistributorBlocks(ctx context.Context, req *disProto.SelectDistributorBlocksRequest) (*disProto.DistributorBlocksList, error) {
	filters := map[string]any{}
	if req.Filters.Status != nil {
		res, err := requests.DistributorStatus(ctx, *req.Filters.Status)
		if err != nil {
			s.Log(ctx).WithError(err).Errorf("invalid distributor status: %s", *req.Filters.Status)

			return nil, err
		}

		filters["status"] = res
	}
	if req.Filters.InitiatorId != nil {
		initiatorID, err := requests.UserID(ctx, *req.Filters.InitiatorId)
		if err != nil {
			s.Log(ctx).WithError(err).Errorf("invalid initiator ID format: %s", *req.Filters.InitiatorId)

			return nil, err
		}

		filters["initiator_id"] = initiatorID
	}
	if req.Filters.DistributorId != nil {
		distributorID, err := requests.DistributorID(ctx, *req.Filters.DistributorId)
		if err != nil {
			s.Log(ctx).WithError(err).Errorf("invalid distributor ID format: %s", *req.Filters.DistributorId)

			return nil, err
		}

		filters["distributor_id"] = distributorID
	}

	blocks, pag, err := s.app.SelectBlockages(ctx, filters, pagi.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to select distributor blocks")

		return nil, err
	}

	return responses.BlocksList(blocks, pag), err
}
