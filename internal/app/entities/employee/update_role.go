package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (e Employee) UpdateEmployeeRole(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	newRole string,
) (models.Employee, error) {
	initiator, err := e.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.Employee{}, err
	}
	user, err := e.GetByUserID(ctx, userID)
	if err != nil {
		return models.Employee{}, err
	}

	if initiator.DistributorID != user.DistributorID {
		return models.Employee{}, errx.ErrorInitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator %s and chosen employee %s have different distributors", initiatorID, userID),
		)
	}

	allowed, err := enum.CompareEmployeeRoles(initiator.Role, user.Role)
	if err != nil {
		return models.Employee{}, errx.EmployeeInvalidRole.Raise(
			fmt.Errorf("new role is invalid: %w", err),
		)
	}
	if allowed != 1 {
		return models.Employee{}, errx.ErrorInitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	allowed, err = enum.CompareEmployeeRoles(initiator.Role, newRole)
	if err != nil {
		return models.Employee{}, errx.EmployeeInvalidRole.Raise(
			fmt.Errorf("new role is invalid: %w", err),
		)
	}
	if allowed != 1 {
		return models.Employee{}, errx.ErrorInitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	err = e.employee.New().FilterUserID(userID).Update(ctx, map[string]interface{}{
		"role":       newRole,
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Employee{}, errx.ErrorEmployeeNotFound.Raise(
				fmt.Errorf("employee with userID %s not found: %w", userID, err),
			)
		default:
			return models.Employee{}, errx.ErrorInternal.Raise(
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
