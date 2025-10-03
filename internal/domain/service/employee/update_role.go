package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) UpdateEmployeeRole(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	newRole string,
) (models.Employee, error) {
	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.Employee{}, err
	}

	user, err := s.Get(ctx, GetFilters{
		UserID: &userID,
	})
	if err != nil {
		return models.Employee{}, err
	}

	if initiator.DistributorID != user.DistributorID {
		return models.Employee{}, errx.ErrorInitiatorIsNotEmployeeOfThisDistributor.Raise(
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
		return models.Employee{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
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
		return models.Employee{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	now := time.Now().UTC()

	err = s.db.UpdateEmployeeRole(ctx, userID, newRole, now)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	return models.Employee{
		UserID:        user.UserID,
		DistributorID: user.DistributorID,
		Role:          newRole,
		UpdatedAt:     now,
		CreatedAt:     user.CreatedAt,
	}, nil
}
