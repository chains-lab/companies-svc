package block

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Crete(
	ctx context.Context,
	initiatorID uuid.UUID,
	companyID uuid.UUID,
	reason string,
) (models.CompanyBlock, error) {
	company, err := s.getCompany(ctx, companyID)
	if err != nil {
		return models.CompanyBlock{}, err
	}

	now := time.Now().UTC()
	block := models.CompanyBlock{
		ID:          uuid.New(),
		CompanyID:   companyID,
		InitiatorID: initiatorID,
		Reason:      reason,
		Status:      enum.CompanyBlockStatusActive,
		BlockedAt:   now,
	}

	activeBlock, err := s.db.GetActiveCompanyBlock(ctx, block.ID)
	if err != nil {
		return models.CompanyBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error, cause: %w", err),
		)
	}
	if !activeBlock.IsNil() {
		return models.CompanyBlock{}, errx.ErrorCompanyHaveAlreadyActiveBlock.Raise(
			fmt.Errorf("company %s already has an active block", companyID),
		)
	}

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.UpdateCompaniesStatus(ctx, companyID, enum.CompanyStatusBlocked, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to updating company status, cause: %w", err),
			)
		}

		err = s.db.CreateCompanyBlock(ctx, block)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed inserting new block, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.CompanyBlock{}, err
	}

	employees, err := s.db.GetCompanyEmployees(ctx, companyID)
	if err != nil {
		return models.CompanyBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company employees for event recipients, cause: %w", err),
		)
	}

	var recipients []uuid.UUID
	for _, emp := range employees.Data {
		recipients = append(recipients, emp.UserID)
	}

	company.Status = enum.CompanyStatusBlocked
	company.UpdatedAt = now

	if err = s.event.PublishCompanyBlocked(ctx, block, company, recipients...); err != nil {
		return models.CompanyBlock{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish company blocked event, cause: %w", err),
		)
	}

	return block, nil
}
