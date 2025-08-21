package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/distributors-svc/pkg/pagination"
	"github.com/google/uuid"
)

type distributorsQ interface {
	New() dbx.DistributorsQ
	Insert(ctx context.Context, input dbx.Distributor) error
	Get(ctx context.Context) (dbx.Distributor, error)
	Select(ctx context.Context) ([]dbx.Distributor, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.DistributorsQ
	FilterStatus(status string) dbx.DistributorsQ
	LikeName(name string) dbx.DistributorsQ

	Page(limit, offset uint64) dbx.DistributorsQ
	Count(ctx context.Context) (uint64, error)

	Transaction(fn func(ctx context.Context) error) error
}

func (a App) GetDistributor(ctx context.Context, ID uuid.UUID) (models.Distributor, error) {
	distributor, err := a.distributor.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(
				ctx,
				fmt.Errorf("distributor %s not found: %w", ID, err),
				ID,
			)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, fmt.Errorf("getting distributor %s: %w", ID, err))
		}
	}

	return models.Distributor{
		ID:        distributor.ID,
		Icon:      distributor.Icon,
		Name:      distributor.Name,
		Status:    distributor.Status,
		UpdatedAt: distributor.UpdatedAt,
		CreatedAt: distributor.CreatedAt,
	}, nil
}

func (a App) UpdateDistributorName(ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	name string,
) (models.Distributor, error) {
	_, err := a.CompareEmployeesRole(ctx, initiatorID, distributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Distributor{}, err
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, map[string]any{
		"name":       name,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return models.Distributor{}, errx.RaiseInternal(ctx, fmt.Errorf("updating distributor name: %w", err))
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(
				ctx,
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
				distributorID,
			)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, fmt.Errorf("getting distributor %s: %w", distributorID, err))
		}
	}

	return models.Distributor{
		ID:        distributor.ID,
		Name:      distributor.Name,
		Icon:      distributor.Icon,
		Status:    distributor.Status,
		UpdatedAt: distributor.UpdatedAt,
		CreatedAt: distributor.CreatedAt,
	}, nil
}

func (a App) UpdateDistributorIcon(ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	icon string,
) (models.Distributor, error) {
	_, err := a.CompareEmployeesRole(ctx, initiatorID, distributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Distributor{}, err
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, map[string]any{
		"icon":       icon,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return models.Distributor{}, errx.RaiseInternal(ctx, fmt.Errorf("updating distributor icon: %w", err))
	}

	return a.GetDistributor(ctx, distributorID)
}

func (a App) SetDistributorStatusInactive(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	_, err := a.CompareEmployeesRole(ctx, initiatorID, distributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Distributor{}, err
	}

	distributor, err := a.distributor.New().FilterID(distributorID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(
				ctx,
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
				distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, fmt.Errorf("getting distributor %s: %w", distributorID, err))
		}
	}
	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.RaiseDistributorStatusBlocked(
			ctx,
			fmt.Errorf("distributor %s is blocked", distributorID),
			distributorID,
		)
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusInactive,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return models.Distributor{}, errx.RaiseInternal(ctx, fmt.Errorf("updating distributor status: %w", err))
	}

	return a.GetDistributor(ctx, distributorID)
}

func (a App) SetDistributorStatusActive(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	_, err := a.CompareEmployeesRole(ctx, initiatorID, distributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Distributor{}, err
	}

	distributor, err := a.distributor.New().FilterID(distributorID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(
				ctx,
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
				distributorID,
			)
		default:
			return models.Distributor{},
				errx.RaiseInternal(ctx, fmt.Errorf("getting distributor %s: %w", distributorID, err))
		}
	}
	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.RaiseDistributorStatusBlocked(
			ctx,
			fmt.Errorf("distributor %s is block", distributorID),
			distributorID,
		)
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusActive,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return models.Distributor{}, errx.RaiseInternal(ctx, fmt.Errorf("updating distributor status: %w", err))
	}

	distributor, err = distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(
				ctx,
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
				distributorID,
			)
		default:
			return models.Distributor{},
				errx.RaiseInternal(ctx, fmt.Errorf("getting distributor %s: %w", distributorID, err))
		}
	}

	return models.Distributor{
		ID:        distributor.ID,
		Name:      distributor.Name,
		Icon:      distributor.Icon,
		Status:    distributor.Status,
		UpdatedAt: distributor.UpdatedAt,
		CreatedAt: distributor.CreatedAt,
	}, nil
}

func (a App) SelectedDistributors(
	ctx context.Context,
	filters map[string]interface{},
	pag pagination.Request,
) ([]models.Distributor, pagination.Response, error) {
	query := a.distributor.New()
	if id, ok := filters["id"]; ok {
		query = query.FilterID(id.(uuid.UUID))
	}
	if status, ok := filters["status"]; ok {
		query = query.FilterStatus(status.(string))
	}
	if name, ok := filters["name"]; ok {
		query = query.LikeName(name.(string))
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	query = query.Page(limit, offset)

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, fmt.Errorf("counting distributors: %w", err))
	}

	distributors, err := query.Select(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, fmt.Errorf("selecting distributors: %w", err))
	}

	var result []models.Distributor
	for _, d := range distributors {
		result = append(result, models.Distributor{
			ID:        d.ID,
			Icon:      d.Icon,
			Name:      d.Name,
			Status:    d.Status,
			UpdatedAt: d.UpdatedAt,
			CreatedAt: d.CreatedAt,
		})
	}

	return result, pagination.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}
