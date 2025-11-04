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

	if company.Status == enum.CompanyStatusBlocked {
		return models.Company{}, errx.ErrorCompanyIsBlocked.Raise(
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

func (s Service) UpdateStatus(
	ctx context.Context,
	companyID uuid.UUID,
	status string,
) (models.Company, error) {
	err := enum.CheckCompanyStatus(status)
	if err != nil {
		return models.Company{}, errx.ErrorInvalidCompanyBlockStatus.Raise(
			fmt.Errorf("failed invalid status %s, cause: %w", status, err),
		)
	}

	if status == enum.CompanyStatusBlocked {
		return models.Company{}, errx.ErrorCannotSetCompanyStatusBlocked.Raise(
			fmt.Errorf("cannot set status to blocked"),
		)
	}

	company, err := s.Get(ctx, companyID)
	if err != nil {
		return models.Company{}, err
	}

	if company.Status == enum.CompanyStatusBlocked {
		return models.Company{}, errx.ErrorCompanyIsBlocked.Raise(
			fmt.Errorf("company %s is blocked", companyID),
		)
	}

	now := time.Now().UTC()
	err = s.db.UpdateCompaniesStatus(ctx, companyID, status, now)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update company status, cause: %w", err),
		)
	}

	switch status {
	case enum.CompanyStatusActive:
		err = s.event.PublishCompanyActivated(ctx, company)
		if err != nil {
			return models.Company{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish company unblocked event, cause: %w", err),
			)
		}

	case enum.CompanyStatusInactive:
		err = s.event.PublishCompanyDeactivated(ctx, company)
		if err != nil {
			return models.Company{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish company blocked event, cause: %w", err),
			)
		}
	}

	return models.Company{
		ID:        company.ID,
		Name:      company.Name,
		Icon:      company.Icon,
		Status:    status,
		UpdatedAt: now,
		CreatedAt: company.CreatedAt,
	}, nil
}
