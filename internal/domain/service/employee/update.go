package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateEmployeeParams struct {
	Position *string
	Label    *string
	Role     *string
}

func (s Service) UpdateEmployee(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	params UpdateEmployeeParams,
) (models.Employee, error) {
	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.Employee{}, err
	}

	employee, err := s.Get(ctx, GetParams{
		UserID: &userID,
	})
	if err != nil {
		return models.Employee{}, err
	}

	if initiator.CompanyID != employee.CompanyID {
		return models.Employee{}, errx.ErrorInitiatorIsNotEmployeeOfThisCompany.Raise(
			fmt.Errorf("initiator %s and chosen employee %s have different companies", initiatorID, userID),
		)
	}

	access, err := enum.CompareEmployeeRoles(initiator.Role, employee.Role)
	if err != nil {
		return models.Employee{}, errx.ErrorInvalidEmployeeRole.Raise(
			fmt.Errorf("compare employee roles: %w", err),
		)
	}
	if access != 1 {
		return models.Employee{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee"),
		)
	}

	if params.Position == nil {
		employee.Position = params.Position
	}
	if params.Label == nil {
		employee.Label = params.Label
	}
	if params.Role != nil {
		access, err = enum.CompareEmployeeRoles(initiator.Role, *params.Role)
		if err != nil {
			return models.Employee{}, errx.EmployeeInvalidRole.Raise(
				fmt.Errorf("new role is invalid: %w", err),
			)
		}
		if access != 1 {
			return models.Employee{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
				fmt.Errorf("initiator have not enough rights to update employee role"),
			)
		}

		employee.Role = *params.Role
	}

	now := time.Now().UTC()

	err = s.db.UpdateEmployee(ctx, userID, params, now)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update employee role, cause: %w", err),
		)
	}

	if err = s.event.PublishEmployeeUpdated(ctx, employee); err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee be kafka, cause: %w", err),
		)
	}

	employee.UpdatedAt = now

	return employee, nil
}
