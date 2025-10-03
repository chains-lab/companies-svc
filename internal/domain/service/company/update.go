package company

import (
	"context"
	"time"

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

	mow := time.Now().UTC()

	if params.Name != nil {
		company.Name = *params.Name
	}

	if params.Icon != nil {
		company.Icon = *params.Icon
	}

	err = s.db.UpdateCompany(ctx, companyID, params, mow)
	if err != nil {
		return models.Company{}, err
	}

	return company, nil
}
