package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, initiatorID, userID, companyID uuid.UUID) error {
	employee, err := s.Get(ctx, GetParams{
		UserID:    &userID,
		CompanyID: &companyID,
	})
	if err != nil {
		return err
	}

	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return err
	}

	if initiator.UserID == userID {
		return errx.ErrorCannotDeleteYourself.Raise(
			fmt.Errorf("initiator %s is trying to delete himself", initiatorID),
		)
	}

	if initiator.CompanyID != employee.CompanyID {
		return errx.ErrorInitiatorIsNotEmployeeOfThisCompany.Raise(
			fmt.Errorf("initiator %s and chosen employee %s have different companies", initiatorID, userID),
		)
	}

	allowed, err := enum.CompareEmployeeRoles(initiator.Role, employee.Role)
	if err != nil {
		return errx.ErrorInvalidEmployeeRole.Raise(err)
	}
	if allowed != 1 {
		return errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to delete employee"),
		)
	}

	company, err := s.db.GetCompanyByID(ctx, companyID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by ID %s, cause: %w", companyID, err),
		)
	}

	employees, err := s.db.GetCompanyEmployees(ctx, companyID, enum.EmployeeRoleAdmin, enum.EmployeeRoleOwner)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company employees, cause: %w", err),
		)
	}

	var recipientsIDs []uuid.UUID
	for _, emp := range employees.Data {
		recipientsIDs = append(recipientsIDs, emp.UserID)
	}

	if err = s.db.DeleteEmployee(ctx, userID, companyID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete employee, cause: %w", err),
		)
	}

	if err = s.event.PublishEmployeeDeleted(ctx, company, employee, recipientsIDs); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish company deleted event, cause: %w", err),
		)
	}

	return nil
}
