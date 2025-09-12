package distributor

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
)

type FilterList struct {
	Name     *string
	Statuses []string
}

func (d Distributor) List(
	ctx context.Context,
	filters FilterList,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Distributor, pagi.Response, error) {
	query := d.distributor.New()

	if filters.Name != nil {
		query = query.LikeName(*filters.Name)
	}
	if filters.Statuses != nil {
		for _, status := range filters.Statuses {
			err := enum.CheckDistributorStatus(status)
			if err != nil {
				return nil, pagi.Response{}, errx.InvalidDistributorStatus.Raise(err)
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

	query = query.Page(limit, offset)

	for _, sort := range sort {
		ascend := sort.Ascend
		switch sort.Field {
		case "name":
			query = query.OrderByName(ascend)
		default:

		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	result := make([]models.Distributor, 0, len(rows))
	for _, dist := range rows {
		result = append(result, models.Distributor{
			ID:        dist.ID,
			Icon:      dist.Icon,
			Name:      dist.Name,
			Status:    dist.Status,
			UpdatedAt: dist.UpdatedAt,
			CreatedAt: dist.CreatedAt,
		})
	}

	return result, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}
