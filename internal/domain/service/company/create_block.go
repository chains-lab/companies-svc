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

func (s Service) CreteBlock(
	ctx context.Context,
	initiatorID uuid.UUID,
	companyID uuid.UUID,
	reason string,
) (models.CompanyBlock, error) {
	_, err := s.Get(ctx, companyID)
	if err != nil {
		return models.CompanyBlock{}, err
	}

	now := time.Now().UTC()

	block := models.CompanyBlock{
		ID:          uuid.New(),
		CompanyID:   companyID,
		InitiatorID: initiatorID,
		Reason:      reason,
		Status:      enum.DistributorBlockStatusActive,
		BlockedAt:   now,
	}

	activeBlock, err := s.db.GetActiveCompanyBlock(ctx, block.ID)
	if err != nil {
		return models.CompanyBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}
	if !activeBlock.IsNil() {
		return models.CompanyBlock{}, errx.ErrorcompanyHaveAlreadyActiveBlock.Raise(
			fmt.Errorf("company %s already has an active block", companyID),
		)
	}

	trErr := s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.UpdateCompaniesStatus(ctx, companyID, enum.DistributorStatusActive, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating company status: %w", err),
			)
		}

		err = s.db.CreateCompanyBlock(ctx, block)
		if err != nil {
			return errx.ErrorInternal.Raise(fmt.Errorf("inserting new block: %w", err))
		}

		return nil
	})
	if trErr != nil {
		return models.CompanyBlock{}, trErr
	}

	return block, nil
}
