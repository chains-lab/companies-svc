package company

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, companyID uuid.UUID) error {
	company, err := s.Get(ctx, companyID)
	if err != nil {
		return err
	}

	if company.Status != enum.CompanyStatusInactive {
		return errx.ErrorOnlyInactiveCompanyCanBeDeleted.Raise(
			fmt.Errorf("only inactive company can be deleted, current status: %s", company.Status),
		)
	}

	var employees models.EmployeesCollection
	var recipientIDs []uuid.UUID

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		if err = s.db.DeleteCompany(ctx, companyID); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete company, cause: %w", err),
			)
		}

		employees, err = s.db.GetCompanyEmployees(ctx, companyID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get company employees, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return err
	}

	for _, emp := range employees.Data {
		recipientIDs = append(recipientIDs, emp.UserID)
	}

	if err = s.event.PublishCompanyDeleted(ctx, company, recipientIDs); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish company deleted event, cause: %w", err),
		)
	}

	return nil
}
