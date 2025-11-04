package company

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
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

	err = s.db.DeleteCompany(ctx, companyID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete company, cause: %w", err),
		)
	}

	if err = s.event.PublishCompanyDeleted(ctx, company); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish company deleted event, cause: %w", err),
		)
	}

	return nil
}
