package handlers

import (
	"context"

	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/pkg/pagination"
)

func (s Service) SelectDistributors(ctx context.Context, req *disProto.SelectDistributorsRequest) (*disProto.DistributorsList, error) {
	filters := map[string]any{}
	if req.Filters.Status != nil {
		filters["status"] = req.Filters.Status
	}
	if req.Filters.NameLike != nil {
		filters["name"] = req.Filters.NameLike
	}

	distributors, pag, err := s.app.SelectDistributors(ctx, filters, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to select distributors")

		return nil, err
	}

	return responses.DistributorsList(distributors, pag), nil
}
