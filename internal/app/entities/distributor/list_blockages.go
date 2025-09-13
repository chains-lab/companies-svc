package distributor

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

type FilterBlockagesList struct {
	Distributors []uuid.UUID
	Initiators   []uuid.UUID
	Statuses     []string
}

func (d Distributor) ListBlockages(
	ctx context.Context,
	filters FilterBlockagesList,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Block, pagi.Response, error) {
	query := d.block.New()

	if filters.Distributors != nil {
		query = query.FilterDistributorID(filters.Distributors...)
	}
	if filters.Initiators != nil {
		query = query.FilterInitiatorID(filters.Initiators...)
	}
	if filters.Statuses != nil {
		for _, status := range filters.Statuses {
			err := enum.CheckDistributorBlockStatus(status)
			if err != nil {
				return nil, pagi.Response{}, errx.ErrorInvalidDistributorBlockStatus.Raise(err)
			}
		}
		query = query.FilterStatus(filters.Statuses...)
	}

	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	for _, s := range sort {
		ascend := s.Ascend
		switch s.Field {
		case "blocked_at":
			query = query.OrderByBlockedAt(ascend)
		default:

		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("counting rows: %w", err),
		)
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("selecting rows: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	res := make([]models.Block, 0, len(rows))
	for _, s := range rows {
		el := models.Block{
			ID:            s.ID,
			DistributorID: s.DistributorID,
			InitiatorID:   s.InitiatorID,
			Reason:        s.Reason,
			Status:        s.Status,
			BlockedAt:     s.BlockedAt,
		}
		if s.CanceledAt != nil {
			canceledAt := *s.CanceledAt
			el.CanceledAt = &canceledAt
		}

		res = append(res, el)
	}

	return res, pagi.Response{
		Total: count,
		Page:  pag.Page,
		Size:  limit,
	}, nil
}
