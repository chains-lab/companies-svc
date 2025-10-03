package company

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

func (s Service) CancelBlock(
	ctx context.Context,
	companyID uuid.UUID,
) (models.Company, error) {
	dis, err := s.Get(ctx, companyID)
	if err != nil {
		return models.Company{}, err
	}

	_, err = s.GetActivecompanyBlock(ctx, companyID)
	if err != nil {
		return models.Company{}, err
	}

	canceledAt := time.Now().UTC()

	trErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.CancelActiveCompanyBlock(ctx, companyID, canceledAt)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("cancelling active company block: %w", err),
			)
		}

		err = s.db.UpdateCompaniesStatus(ctx, companyID, enum.DistributorStatusInactive, canceledAt)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating company status: %w", err),
			)
		}

		return nil
	})
	if trErr != nil {
		return models.Company{}, trErr
	}

	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("updating company status: %w", err),
		)
	}

	dis.Status = enum.DistributorStatusInactive
	dis.UpdatedAt = canceledAt

	return dis, nil
}
