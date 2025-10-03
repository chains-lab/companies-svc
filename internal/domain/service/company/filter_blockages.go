package company

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type FilterBlockages struct {
	CompanyID   *uuid.UUID
	InitiatorID *uuid.UUID
	Status      *string
}

func (s Service) FilterBlockages(
	ctx context.Context,
	filters FilterBlockages,
	page, size uint64,
) (models.CompanyBlockCollection, error) {
	res, err := s.db.FilterCompanyBlocks(ctx, filters, page, size)
	if err != nil {
		return models.CompanyBlockCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed filter company blocks, cause: %w", err),
		)
	}

	return res, err
}
