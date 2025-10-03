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

func (s Service) UpdateStatus(
	ctx context.Context,
	companyID uuid.UUID,
	status string,
) (models.Company, error) {
	err := enum.CheckDistributorStatus(status)
	if err != nil {
		return models.Company{}, errx.ErrorInvalidcompanyBlockStatus.Raise(
			fmt.Errorf("invalid status %s: %w", status, err),
		)
	}

	if status == enum.DistributorStatusBlocked {
		return models.Company{}, errx.ErrorCannotSetcompaniestatusBlocked.Raise(
			fmt.Errorf("cannot set status to blocked"),
		)
	}

	company, err := s.Get(ctx, companyID)
	if err != nil {
		return models.Company{}, err
	}

	if company.Status == enum.DistributorStatusBlocked {
		return models.Company{}, errx.ErrorcompanyIsBlocked.Raise(
			fmt.Errorf("company %s is blocked", companyID),
		)
	}

	now := time.Now().UTC()
	err = s.db.UpdateCompaniesStatus(ctx, companyID, status, now)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	return models.Company{
		ID:        company.ID,
		Name:      company.Name,
		Icon:      company.Icon,
		Status:    enum.DistributorStatusInactive,
		UpdatedAt: now,
		CreatedAt: company.CreatedAt,
	}, nil
}
