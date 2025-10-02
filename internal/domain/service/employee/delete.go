package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) Delete(
	ctx context.Context,
	initiatorID, userID, distributorID uuid.UUID,
) error {
	if initiatorID == userID {
		return errx.ErrorCannotDeleteYourself.Raise(
			fmt.Errorf("initiatorID and userID are the same: %s", initiatorID),
		)
	}

	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return err
	}
	user, err := s.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if initiator.DistributorID != user.DistributorID || initiator.DistributorID != distributorID {
		return errx.ErrorInitiatorAndUserHaveDifferentDistributors.Raise(
			fmt.Errorf("employee with userID %s not found in distributor %s", userID, initiator.DistributorID),
		)
	}

	allowed, err := enum.CompareEmployeeRoles(initiator.Role, user.Role)
	if err != nil {
		return errx.EmployeeInvalidRole.Raise(
			fmt.Errorf("new role is invalid: %w", err),
		)
	}
	if allowed != 1 {
		return errx.ErrorInitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	err = s.employee.New().FilterUserID(userID).FilterDistributorID(initiator.DistributorID).Delete(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errx.ErrorEmployeeNotFound.Raise(
				fmt.Errorf("employee with userID %s not found in distributor %s: %w", userID, initiator.DistributorID, err),
			)
		default:
			return errx.ErrorInternal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}

	return nil
}

type DeleteManyParams struct {
	DistributorID *uuid.UUID
	Roles         []string
}

func (s Service) DeleteMany(ctx context.Context, params DeleteManyParams) error {
	q := s.employee.New()
	if params.DistributorID != nil {
		q = q.FilterDistributorID(*params.DistributorID)
	}
	if params.Roles != nil && len(params.Roles) > 0 {
		q = q.FilterRole(params.Roles...)
	}

	err := q.Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("deleting employees: %w", err),
		)
	}
	return nil
}
