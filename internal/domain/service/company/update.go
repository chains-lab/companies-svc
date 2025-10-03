package company

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type UpdateParams struct {
	Name *string
	Icon *string
}

func (s Service) Update(ctx context.Context,
	companyID uuid.UUID,
	params UpdateParams,
) (models.Company, error) {
	company, err := s.Get(ctx, companyID)
	if err != nil {
		return models.Company{}, err
	}

	if company.Status == enum.DistributorStatusBlocked {
		return models.Company{}, errx.ErrorcompanyIsBlocked.Raise(
			fmt.Errorf("company with ID %s is blocked", companyID),
		)
	}

	mow := time.Now().UTC()

	if params.Name != nil {
		company.Name = *params.Name
	}

	if params.Icon != nil {
		company.Icon = *params.Icon
	}

	err = s.db.UpdateCompany(ctx, companyID, params, mow)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update company status, cause: %w", err),
		)
	}

	return company, nil
}
