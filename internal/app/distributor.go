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
	"github.com/chains-lab/pagi"
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
	FilterStatus(status ...string) dbx.DistributorsQ
	LikeName(name string) dbx.DistributorsQ

	OrderByName(ascend bool) dbx.DistributorsQ

	Page(limit, offset uint64) dbx.DistributorsQ
	Count(ctx context.Context) (uint64, error)

	Transaction(fn func(ctx context.Context) error) error
}

func (a App) CreateDistributor(
	ctx context.Context,
	initiatorID uuid.UUID,
	name string,
	icon string,
) (models.Distributor, error) {
	_, err := a.employee.New().FilterUserID(initiatorID).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Distributor{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}
	if err == nil {
		return models.Distributor{}, errx.EmployeeAlreadyExists.Raise(
			fmt.Errorf("employee with userID %s already exists", initiatorID),
		)
	}

	stmt := dbx.Distributor{
		ID:        uuid.New(),
		Name:      name,
		Icon:      icon,
		Status:    enum.DistributorStatusActive,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = a.distributor.New().Insert(ctx, stmt)
	if err != nil {
		return models.Distributor{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	return models.Distributor{
		ID:        stmt.ID,
		Name:      stmt.Name,
		Icon:      stmt.Icon,
		Status:    stmt.Status,
		CreatedAt: stmt.CreatedAt,
		UpdatedAt: stmt.UpdatedAt,
	}, err
}

func (a App) GetDistributor(ctx context.Context, ID uuid.UUID) (models.Distributor, error) {
	distributor, err := a.distributor.New().FilterID(ID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.DistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", ID, err),
			)
		default:
			return models.Distributor{}, errx.Internal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
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

type SelectDistributorsParams struct {
	Name     *string
	Statuses []string
}

func (a App) SelectDistributors(
	ctx context.Context,
	filters SelectDistributorsParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Distributor, pagi.Response, error) {
	query := a.distributor.New()

	if filters.Name != nil {
		query = query.LikeName(*filters.Name)
	}
	if filters.Statuses != nil {
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
		return nil, pagi.Response{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	result := make([]models.Distributor, 0, len(rows))
	for _, d := range rows {
		result = append(result, models.Distributor{
			ID:        d.ID,
			Icon:      d.Icon,
			Name:      d.Name,
			Status:    d.Status,
			UpdatedAt: d.UpdatedAt,
			CreatedAt: d.CreatedAt,
		})
	}

	return result, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}

type UpdateDistributorInput struct {
	Name *string
	Icon *string
}

func (a App) UpdateDistributor(ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	input UpdateDistributorInput,
) (models.Distributor, error) {
	update := map[string]any{}

	if input.Name != nil {
		update["name"] = *input.Name
	}
	if input.Icon != nil {
		update["icon"] = *input.Icon
	}
	update["updated_at"] = time.Now().UTC()

	initiator, err := a.GetInitiatorEmployee(ctx, initiatorID)
	if err != nil {
		return models.Distributor{}, err
	}

	access, err := enum.ComparisonEmployeeRoles(initiator.Role, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Distributor{}, errx.EmployeeRoleNotSupported.Raise(
			fmt.Errorf("initiator %s have invalid role %s: %w", initiatorID, initiator.Role, err),
		)
	}
	if access < 0 {
		return models.Distributor{}, errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator %s have not enough permissions to update distributor %s", initiatorID, distributorID),
		)
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, update)
	if err != nil {
		return models.Distributor{}, errx.Internal.Raise(fmt.Errorf("updating distributor name: %w", err))
	}

	distributor, err := distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.DistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
			)
		default:
			return models.Distributor{}, errx.Internal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
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
	initiator, err := a.GetInitiatorEmployee(ctx, initiatorID)
	if err != nil {
		return models.Distributor{}, err
	}

	if initiator.Role != enum.EmployeeRoleOwner {
		return models.Distributor{}, errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator %s have not enough permissions to update distributor %s", initiatorID, distributorID),
		)
	}

	distributor, err := a.distributor.New().FilterID(distributorID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.DistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
			)
		default:
			return models.Distributor{}, errx.Internal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}
	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.DistributorStatusBlocked.Raise(
			fmt.Errorf("distributor %s is blocked", distributorID),
		)
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusInactive,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return models.Distributor{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	return a.GetDistributor(ctx, distributorID)
}

func (a App) SetDistributorStatusActive(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
) (models.Distributor, error) {
	initiator, err := a.GetInitiatorEmployee(ctx, initiatorID)
	if err != nil {
		return models.Distributor{}, err
	}

	if initiator.Role != enum.EmployeeRoleOwner {
		return models.Distributor{}, errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator %s have not enough permissions to update distributor %s", initiatorID, distributorID),
		)
	}

	distributor, err := a.distributor.New().FilterID(distributorID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.DistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
			)
		default:
			return models.Distributor{},
				errx.Internal.Raise(
					fmt.Errorf("internal error: %w", err),
				)
		}
	}
	if distributor.Status == enum.DistributorStatusBlocked {
		return models.Distributor{}, errx.DistributorStatusBlocked.Raise(
			fmt.Errorf("distributor %s is block", distributorID),
		)
	}

	distributorQ := a.distributor.New().FilterID(distributorID)

	err = distributorQ.Update(ctx, map[string]any{
		"status":     enum.DistributorStatusActive,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return models.Distributor{}, errx.Internal.Raise(fmt.Errorf("updating distributor status: %w", err))
	}

	distributor, err = distributorQ.Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Distributor{}, errx.DistributorNotFound.Raise(
				fmt.Errorf("distributor %s not found: %w", distributorID, err),
			)
		default:
			return models.Distributor{},
				errx.Internal.Raise(
					fmt.Errorf("internal error: %w", err),
				)
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
