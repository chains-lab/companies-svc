package company

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Get(ctx context.Context, ID uuid.UUID) (models.Company, error) {
	company, err := s.db.GetCompanyByID(ctx, ID)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by EmployeeID, cause: %w", err),
		)
	}

	if company.IsNil() {
		return models.Company{}, errx.ErrorCompanyNotFound.Raise(
			fmt.Errorf("company with EmployeeID %s not found", ID),
		)
	}

	return company, nil
}
