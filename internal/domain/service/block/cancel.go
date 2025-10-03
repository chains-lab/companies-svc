package block

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Cancel(ctx context.Context, companyID uuid.UUID) (models.Company, error) {
	dis, err := s.getCompany(ctx, companyID)
	if err != nil {
		return models.Company{}, err
	}

	_, err = s.GetActiveCompanyBlock(ctx, companyID)
	if err != nil {
		return models.Company{}, err
	}

	now := time.Now().UTC()

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.CancelActiveCompanyBlock(ctx, companyID, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to cancelling active company block, cause: %w", err),
			)
		}

		err = s.db.UpdateCompaniesStatus(ctx, companyID, enum.CompanyStatusInactive, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to updating company status, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Company{}, err
	}

	dis.Status = enum.CompanyStatusInactive
	dis.UpdatedAt = now

	return dis, nil
}
