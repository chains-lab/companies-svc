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

func (s Service) UpdateRole(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	role string,
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

	company, err := s.db.GetCompanyByID(ctx, initiator.CompanyID)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by id %s, cause: %w", initiator.CompanyID, err),
		)
	}
	if company.IsNil() {
		return models.Employee{}, errx.ErrorCompanyNotFound.Raise(
			fmt.Errorf("company with id %s not found", initiator.CompanyID),
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

	access, err = enum.CompareEmployeeRoles(initiator.Role, role)
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

	employee.Role = role
	employee.UpdatedAt = time.Now().UTC()

	if err = s.db.UpdateEmployeeRole(ctx, employee.UserID, role, employee.UpdatedAt); err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update employee role, cause: %w", err),
		)
	}

	employees, err := s.db.GetCompanyEmployees(
		ctx,
		company.ID,
		enum.EmployeeRoleOwner,
		enum.EmployeeRoleAdmin,
	)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company employees by role %s, cause: %w", employee.Role, err),
		)
	}

	var recipientIDs []uuid.UUID
	for _, emp := range employees.Data {
		recipientIDs = append(recipientIDs, emp.UserID)
	}
	recipientIDs = append(recipientIDs, initiatorID)

	if err = s.event.PublishEmployeeUpdated(ctx, company, initiator, recipientIDs); err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee be kafka, cause: %w", err),
		)
	}

	return employee, nil
}
