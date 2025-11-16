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

type UpdateParams struct {
	Position *string
	Label    *string
	Role     *string
}

func (s Service) UpdateByEmployee(
	ctx context.Context,
	userID uuid.UUID,
	initiatorID uuid.UUID,
	params UpdateParams,
) (models.Employee, error) {
	initiator, err := s.validateInitiatorRight(ctx, initiatorID, nil, enum.EmployeeRoleOwner, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Employee{}, err
	}

	employee, err := s.Get(ctx, GetParams{
		UserID:    &userID,
		CompanyID: &initiator.CompanyID,
	})
	if err != nil {
		return models.Employee{}, err
	}

	access, err := enum.CompareEmployeeRoles(initiator.Role, employee.Role)
	if err != nil {
		return models.Employee{}, errx.ErrorInvalidEmployeeRole.Raise(
			fmt.Errorf("compare employee roles: %w", err),
		)
	}
	if access != 1 {
		return models.Employee{}, errx.ErrorNotEnoughRight.Raise(
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
			return models.Employee{}, errx.ErrorInvalidEmployeeRole.Raise(
				fmt.Errorf("compare employee roles: %w", err),
			)
		}
		if access != 1 {
			return models.Employee{}, errx.ErrorNotEnoughRight.Raise(
				fmt.Errorf("initiator have not enough rights to update employee role"),
			)
		}

		employee.Role = *params.Role
	}
	updatedAt := time.Now().UTC()
	employee.UpdatedAt = updatedAt

	company, err := s.getCompany(ctx, initiator.CompanyID)
	if err != nil {
		return models.Employee{}, err
	}
	if company.Status != enum.CompanyStatusActive {
		return models.Employee{}, errx.ErrorCompanyIsNotActive.Raise(
			fmt.Errorf("company with id %s is not active", initiator.CompanyID),
		)
	}

	return s.update(ctx, employee, company, params, initiatorID)
}

type UpdateMyParams struct {
	Position *string
	Label    *string
}

func (s Service) UpdateMy(
	ctx context.Context,
	initiatorID uuid.UUID,
	params UpdateMyParams,
) (models.Employee, error) {
	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.Employee{}, err
	}

	if params.Label != nil {
		initiator.Label = params.Label
	}
	if params.Position != nil {
		initiator.Position = params.Position
	}
	updatedAt := time.Now().UTC()
	initiator.UpdatedAt = updatedAt

	company, err := s.getCompany(ctx, initiator.CompanyID)
	if err != nil {
		return models.Employee{}, err
	}
	//if company.Status != enum.CompanyStatusActive {
	//	return models.Employee{}, errx.ErrorCompanyIsNotActive.Raise(
	//		fmt.Errorf("company with id %s is not active", initiator.CompanyID),
	//	)
	//}

	return s.update(ctx, initiator, company, UpdateParams{
		Position: params.Position,
		Label:    params.Label,
	}, initiatorID)
}

func (s Service) update(
	ctx context.Context,
	employee models.Employee,
	company models.Company,
	params UpdateParams,
	recipientIDs ...uuid.UUID,
) (models.Employee, error) {
	if params.Position == nil {
		employee.Position = params.Position
	}
	if params.Label == nil {
		employee.Label = params.Label
	}
	updatedAt := time.Now().UTC()
	employee.UpdatedAt = updatedAt

	if err := s.db.UpdateEmployee(ctx, employee.UserID, params, updatedAt); err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update employee role, cause: %w", err),
		)
	}
	if err := s.event.PublishEmployeeUpdated(ctx, company, employee, recipientIDs...); err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee be kafka, cause: %w", err),
		)
	}

	return employee, nil
}
