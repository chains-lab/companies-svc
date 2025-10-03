package company

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetBlock(
	ctx context.Context,
	BlockID uuid.UUID,
) (models.CompanyBlock, error) {
	block, err := s.db.GetCompanyBlockByID(ctx, BlockID)
	if err != nil {
		return models.CompanyBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting block with ID %s: %w", BlockID, err),
		)
	}

	if block.IsNil() {
		return models.CompanyBlock{}, errx.ErrorcompanyBlockNotFound.Raise(
			fmt.Errorf("block with ID %s not found", BlockID),
		)
	}

	return block, nil
}

func (s Service) GetActivecompanyBlock(ctx context.Context, companyID uuid.UUID) (models.CompanyBlock, error) {
	block, err := s.db.GetActiveCompanyBlock(ctx, companyID)
	if err != nil {
		return models.CompanyBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting active block for company with ID %s, cause: %w", companyID, err),
		)
	}

	if block.IsNil() {
		return models.CompanyBlock{}, errx.ErrorcompanyBlockNotFound.Raise(
			fmt.Errorf("active block for company with ID %s not found", companyID),
		)
	}

	return block, nil
}
