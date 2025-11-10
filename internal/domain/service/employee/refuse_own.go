package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) RefuseMe(ctx context.Context, initiatorID uuid.UUID) error {
	own, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return err
	}

	if own.Role == enum.EmployeeRoleOwner {
		return errx.ErrorOwnerCannotRefuseSelf.Raise(
			fmt.Errorf("owner cannot refuse self"),
		)
	}

	company, err := s.db.GetCompanyByID(ctx, own.CompanyID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by ID %s, cause: %w", own.CompanyID, err),
		)
	}

	employees, err := s.db.GetCompanyEmployees(ctx, own.CompanyID, enum.EmployeeRoleAdmin, enum.EmployeeRoleOwner)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company employees, cause: %w", err),
		)
	}

	var recipientsIDs []uuid.UUID
	for _, emp := range employees.Data {
		recipientsIDs = append(recipientsIDs, emp.UserID)
	}

	err = s.db.DeleteEmployee(ctx, own.UserID, own.CompanyID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee in transaction, cause: %w", err),
		)
	}

	if err = s.event.PublishEmployeeDeleted(ctx, company, own, recipientsIDs); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee be kafka, cause: %w", err),
		)
	}

	return nil
}
