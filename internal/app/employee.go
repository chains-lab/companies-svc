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

type employeesQ interface {
	New() dbx.EmployeeQ

	Insert(ctx context.Context, input dbx.Employee) error
	Get(ctx context.Context) (dbx.Employee, error)
	Select(ctx context.Context) ([]dbx.Employee, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterUserID(id uuid.UUID) dbx.EmployeeQ
	FilterDistributorID(id ...uuid.UUID) dbx.EmployeeQ
	FilterRole(role ...string) dbx.EmployeeQ

	OrderByRole(ascend bool) dbx.EmployeeQ

	Page(limit, offset uint64) dbx.EmployeeQ
	Count(ctx context.Context) (uint64, error)

	Transaction(fn func(ctx context.Context) error) error
}

func (a App) GetInitiatorEmployee(
	ctx context.Context,
	initiatorID uuid.UUID,
) (models.Employee, error) {
	employee, err := a.employee.New().
		FilterUserID(initiatorID).
		Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Employee{}, errx.InitiatorNotEmployee.Raise(
				fmt.Errorf("initiator with userID %s not found: %w", initiatorID, err),
			)
		default:
			return models.Employee{}, errx.Internal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}

	return models.Employee{
		UserID:        employee.UserID,
		DistributorID: employee.DistributorID,
		Role:          employee.Role,
		UpdatedAt:     employee.UpdatedAt,
		CreatedAt:     employee.CreatedAt,
	}, nil
}

func (a App) GetEmployee(ctx context.Context, userID uuid.UUID) (models.Employee, error) {
	employee, err := a.employee.New().
		FilterUserID(userID).
		Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Employee{}, errx.EmployeeNotFound.Raise(
				fmt.Errorf("employee with userID %s not found: %w", userID, err),
			)
		default:
			return models.Employee{}, errx.Internal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}

	return models.Employee{
		UserID:        employee.UserID,
		DistributorID: employee.DistributorID,
		Role:          employee.Role,
		UpdatedAt:     employee.UpdatedAt,
		CreatedAt:     employee.CreatedAt,
	}, nil
}

type SelectEmployeesParams struct {
	Distributors []uuid.UUID
	Roles        []string
}

func (a App) SelectEmployees(
	ctx context.Context,
	filters SelectEmployeesParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.Employee, pagi.Response, error) {
	query := a.employee.New()

	if filters.Distributors != nil {
		query = query.FilterDistributorID(filters.Distributors...)
	}
	if filters.Roles != nil {
		query = query.FilterRole(filters.Roles...)
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

	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	query = query.Page(limit, offset)

	for _, sort := range sort {
		ascend := sort.Ascend
		switch sort.Field {
		case "role":
			query = query.OrderByRole(ascend)
		default:

		}
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

	var result []models.Employee
	for _, emp := range rows {
		result = append(result, models.Employee{
			UserID:        emp.UserID,
			DistributorID: emp.DistributorID,
			Role:          emp.Role,
			UpdatedAt:     emp.UpdatedAt,
			CreatedAt:     emp.CreatedAt,
		})
	}

	return result, pagi.Response{
		Total: total,
		Page:  pag.Page,
		Size:  pag.Size,
	}, nil
}

func (a App) UpdateEmployeeRole(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	newRole string,
) (models.Employee, error) {
	initiator, err := a.GetInitiatorEmployee(ctx, initiatorID)
	if err != nil {
		return models.Employee{}, err
	}
	user, err := a.GetEmployee(ctx, userID)
	if err != nil {
		return models.Employee{}, err
	}

	if initiator.DistributorID != user.DistributorID {
		return models.Employee{}, errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator %s and chosen employee %s have different distributors", initiatorID, userID),
		)
	}

	allowed, err := enum.ComparisonEmployeeRoles(initiator.Role, user.Role)
	if err != nil {
		return models.Employee{}, errx.EmployeeRoleNotSupported.Raise(
			fmt.Errorf("new role is invalid: %w", err),
		)
	}
	if allowed != 1 {
		return models.Employee{}, errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	allowed, err = enum.ComparisonEmployeeRoles(initiator.Role, newRole)
	if err != nil {
		return models.Employee{}, errx.EmployeeRoleNotSupported.Raise(
			fmt.Errorf("new role is invalid: %w", err),
		)
	}
	if allowed != 1 {
		return models.Employee{}, errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	err = a.employee.New().FilterUserID(userID).Update(ctx, map[string]interface{}{
		"role":       newRole,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Employee{}, errx.EmployeeNotFound.Raise(
				fmt.Errorf("employee with userID %s not found: %w", userID, err),
			)
		default:
			return models.Employee{}, errx.Internal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}

	return models.Employee{
		UserID:        user.UserID,
		DistributorID: user.DistributorID,
		Role:          newRole,
		UpdatedAt:     time.Now().UTC(),
		CreatedAt:     user.CreatedAt,
	}, nil
}

func (a App) DeleteEmployee(ctx context.Context, initiatorID, userID, distributorID uuid.UUID) error {
	initiator, err := a.GetInitiatorEmployee(ctx, initiatorID)
	if err != nil {
		return err
	}
	user, err := a.GetEmployee(ctx, userID)
	if err != nil {
		return err
	}

	allowed, err := enum.ComparisonEmployeeRoles(initiator.Role, user.Role)
	if err != nil {
		return errx.EmployeeRoleNotSupported.Raise(
			fmt.Errorf("new role is invalid: %w", err),
		)
	}
	if allowed != 1 {
		return errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	err = a.employee.New().FilterUserID(userID).FilterDistributorID(distributorID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.EmployeeNotFound.Raise(
				fmt.Errorf("employee with userID %s not found in distributor %s: %w", userID, distributorID, err),
			)
		default:
			return errx.Internal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}

	return nil
}
