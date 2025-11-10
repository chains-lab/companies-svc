package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) UpdateOwn(
	ctx context.Context,
	initiatorID uuid.UUID,
	params UpdateEmployeeParams,
) (models.Employee, error) {
	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.Employee{}, err
	}

	if params.Label != nil {
		initiator.Label = params.Label
	}
	if params.Position != nil {
		initiator.Position = params.Position
	}
	updatedAt := time.Now().UTC()
	initiator.UpdatedAt = updatedAt

	company, err := s.db.GetCompanyByID(ctx, initiator.CompanyID)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by id %s, cause: %w", initiator.CompanyID, err),
		)
	}
	if company.IsNil() {
		return models.Employee{}, errx.ErrorCompanyNotFound.Raise(
			fmt.Errorf("company with id %s not found", initiator.CompanyID),
		)
	}

	err = s.db.UpdateEmployee(ctx, initiatorID, params, updatedAt)
	if err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update employee %s: %w", initiatorID, err),
		)
	}

	if err = s.event.PublishEmployeeUpdated(ctx, company, initiator, []uuid.UUID{initiatorID}); err != nil {
		return models.Employee{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee be kafka, cause: %w", err),
		)
	}

	return initiator, nil
}
