package app

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
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
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, ID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
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

func (a App) CreateDistributor(
	ctx context.Context,
	initiatorID uuid.UUID,
	name string,
) (models.Distributor, error) {
	distributor := dbx.Distributor{
		ID:        uuid.New(),
		Icon:      "",
		Name:      name,
		Status:    enum.DistributorStatusActive,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	trErr := a.distributor.Transaction(func(ctx context.Context) error {
		err := a.distributor.New().Insert(ctx, distributor)
		if err != nil {
			//don't need to check for sql.ErrNoRows here, because Insert will not return it
			return errx.RaiseInternal(ctx, err)
		}

		err = a.employee.New().Insert(ctx, dbx.Employee{
			UserID:        initiatorID,
			DistributorID: distributor.ID,
			Role:          enum.EmployeeRoleAdmin,
			UpdatedAt:     time.Now().UTC(),
			CreatedAt:     time.Now().UTC(),
		})
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return errx.RaiseEmployeCantCreateDisrtibutor(ctx, err, initiatorID)
			default:
				return errx.RaiseInternal(ctx, err)
			}
		}

		return nil
	})
	if trErr != nil {
		return models.Distributor{}, trErr
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
		}
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
		}
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
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

func (a App) SetDistributorStatusInactive(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	_, err := a.CompareEmployeesRole(ctx, initiatorID, distributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Distributor{}, err
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusInactive,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
		}
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
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

func (a App) SetDistributorStatusActive(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	_, err := a.CompareEmployeesRole(ctx, initiatorID, distributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Distributor{}, err
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusActive,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
		}
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
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

// SetDistributorStatusSuspend allows a system admin to suspend a distributor.
func (a App) SetDistributorStatusSuspend(
	ctx context.Context,
	distributorID uuid.UUID,
	reason string,
) (models.Distributor, error) {
	distributorQ := a.distributor.New().FilterID(distributorID)

	err := distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusSuspend,
		"reason":     reason,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
		}
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
		}
	}

	//TODO: Implement logic to send email to admin about suspension with reason

	return models.Distributor{
		ID:        distributor.ID,
		Name:      distributor.Name,
		Icon:      distributor.Icon,
		Status:    distributor.Status,
		UpdatedAt: distributor.UpdatedAt,
		CreatedAt: distributor.CreatedAt,
	}, nil
}

// DeleteDistributorStatusSuspend allows a system admin to lift the suspension of a distributor.
func (a App) DeleteDistributorStatusSuspend(ctx context.Context, distributorID uuid.UUID) (models.Distributor, error) {
	distributorQ := a.distributor.New().FilterID(distributorID)

	err := distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusInactive,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
		}
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.RaiseDistributorNotFound(ctx, err, distributorID)
		default:
			return models.Distributor{}, errx.RaiseInternal(ctx, err)
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
