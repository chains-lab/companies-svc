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

type employeesQ interface {
	New() dbx.EmployeeQ

	Insert(ctx context.Context, input dbx.Employee) error
	Get(ctx context.Context) (dbx.Employee, error)
	Select(ctx context.Context) ([]dbx.Employee, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterUserID(id uuid.UUID) dbx.EmployeeQ
	FilterDistributorID(id uuid.UUID) dbx.EmployeeQ
	FilterRole(role string) dbx.EmployeeQ

	OrderByRole(ascend bool) dbx.EmployeeQ

	Page(limit, offset uint64) dbx.EmployeeQ
	Count(ctx context.Context) (uint64, error)

	Transaction(fn func(ctx context.Context) error) error
}

func (a App) GetEmployee(ctx context.Context, userID uuid.UUID) (models.Employee, error) {
	employee, err := a.employee.New().
		FilterUserID(userID).
		Get(ctx)
	if err != nil {
		return models.Employee{}, err
	}

	return models.Employee{
		UserID:        employee.UserID,
		DistributorID: employee.DistributorID,
		Role:          employee.Role,
		UpdatedAt:     employee.UpdatedAt,
		CreatedAt:     employee.CreatedAt,
	}, nil
}

func (a App) GetDistributorEmployee(ctx context.Context, userID, distributorID uuid.UUID) (models.Employee, error) {
	employee, err := a.employee.New().
		FilterUserID(userID).
		FilterDistributorID(distributorID).
		Get(ctx)
	if err != nil {
		return models.Employee{}, err
	}

	return models.Employee{
		UserID:        employee.UserID,
		DistributorID: employee.DistributorID,
		Role:          employee.Role,
		UpdatedAt:     employee.UpdatedAt,
		CreatedAt:     employee.CreatedAt,
	}, nil
}

func (a App) CompareEmployeesRole(
	ctx context.Context,
	initiatorID uuid.UUID,
	distributorID uuid.UUID,
	role string) (
	// returns:
	models.Employee,
	error,
) {
	employee, err := a.employee.New().
		FilterUserID(initiatorID).
		FilterDistributorID(distributorID).
		Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Employee{}, errx.RaiseInitiatorNotEmployee(ctx, err, initiatorID, distributorID)
		default:
			return models.Employee{}, errx.RaiseInternal(ctx, err)
		}
	}

	access, err := enum.ComparisonEmployeeRoles(employee.Role, role)
	if err != nil {
		return models.Employee{}, errx.RaiseInternal(ctx, err)
	}
	if access < 0 {
		return models.Employee{}, errx.RaiseInitiatorEmployeeHaveNotEnoughPermissions(
			ctx,
			fmt.Errorf("initiator have not enough rights"),
			initiatorID,
			distributorID,
		)
	}

	return models.Employee{
		UserID:        employee.UserID,
		DistributorID: employee.DistributorID,
		Role:          employee.Role,
		UpdatedAt:     employee.UpdatedAt,
		CreatedAt:     employee.CreatedAt,
	}, nil
}

func (a App) AllowedToInteractWithEmployee(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	distributorID uuid.UUID,
	// returns:
) (int, error) {
	initiator, err := a.GetDistributorEmployee(ctx, initiatorID, distributorID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return -1, errx.RaiseInitiatorNotEmployee(ctx, err, initiatorID, distributorID)
		default:
			return -1, errx.RaiseInternal(ctx, err)
		}
	}

	user, err := a.GetDistributorEmployee(ctx, userID, distributorID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return -1, errx.RaiseEmployeeNotFound(ctx, err, userID, distributorID)
		default:
			return -1, errx.RaiseInternal(ctx, err)
		}
	}

	access, err := enum.ComparisonEmployeeRoles(initiator.Role, user.Role)
	if err != nil {
		return -1, errx.RaiseInternal(ctx, err)
	}

	return access, nil
}

func (a App) GetEmployees(
	ctx context.Context,
	filters map[string]any,
	pag pagination.Request) (
	// returns:
	[]models.Employee,
	pagination.Response,
	error,
) {

	query := a.employee.New()
	if distributorID, ok := filters["distributor_id"].(uuid.UUID); ok {
		query = query.FilterDistributorID(distributorID)
	}
	if userID, ok := filters["user_id"].(uuid.UUID); ok {
		query = query.FilterUserID(userID)
	}
	if role, ok := filters["role"].(string); ok {
		query = query.FilterRole(role)
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	employees, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagination.Response{}, err
	}
	total, err := query.Count(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
	}

	var result []models.Employee
	for _, emp := range employees {
		result = append(result, models.Employee{
			UserID:        emp.UserID,
			DistributorID: emp.DistributorID,
			Role:          emp.Role,
			UpdatedAt:     emp.UpdatedAt,
			CreatedAt:     emp.CreatedAt,
		})
	}

	return result, pagination.Response{
		Total: total,
		Page:  pag.Page,
		Size:  limit,
	}, nil
}

func (a App) UpdateEmployeeRole(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	distributorID uuid.UUID,
	newRole string) (
	//returns:
	models.Employee,
	error,
) {
	initiator, err := a.CompareEmployeesRole(ctx, initiatorID, distributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Employee{}, err
	}

	allowed, err := a.AllowedToInteractWithEmployee(ctx, initiatorID, userID, distributorID)
	if err != nil {
		return models.Employee{}, err
	}
	if allowed != 1 {
		return models.Employee{}, errx.RaiseInitiatorEmployeeHaveNotEnoughPermissions(
			ctx,
			fmt.Errorf("initiator have not enough rights to update employee role"),
			initiatorID,
			distributorID,
		)
	}

	user, err := a.GetDistributorEmployee(ctx, userID, distributorID)
	if err != nil {
		return models.Employee{}, err
	}

	access, err := enum.ComparisonEmployeeRoles(initiator.Role, user.Role)
	if err != nil {
		return models.Employee{}, err
	}
	if access < 1 {
		return models.Employee{}, err
	}

	err = a.employee.New().FilterUserID(initiator.UserID).FilterDistributorID(distributorID).Update(ctx, map[string]interface{}{
		"role":       newRole,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		return models.Employee{}, err
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
	allowed, err := a.AllowedToInteractWithEmployee(ctx, initiatorID, userID, distributorID)
	if err != nil {
		return err
	}
	if allowed != 1 {
		return errx.RaiseInitiatorEmployeeHaveNotEnoughPermissions(
			ctx,
			fmt.Errorf("initiator have not enough rights to delete employee"),
			initiatorID,
			distributorID,
		)
	}

	err = a.employee.New().FilterUserID(userID).FilterDistributorID(distributorID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
