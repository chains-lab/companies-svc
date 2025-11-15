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

func (s Service) UpdateByInitiator(
	ctx context.Context,
	initiatorID, companyID uuid.UUID,
	params UpdateParams,
) (models.Company, error) {
	_, err := s.validateInitiatorRight(ctx, initiatorID, companyID, enum.EmployeeRoleOwner, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Company{}, err
	}

	return s.update(ctx, companyID, params)
}

func (s Service) update(ctx context.Context,
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

func (s Service) UpdateStatusByInitiator(
	ctx context.Context,
	initiatorID, companyID uuid.UUID,
	status string,
) (models.Company, error) {
	_, err := s.validateInitiatorRight(ctx, initiatorID, companyID, enum.EmployeeRoleOwner)
	if err != nil {
		return models.Company{}, err
	}

	return s.updateStatus(ctx, companyID, status)
}

func (s Service) updateStatus(
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

	employees, err := s.db.GetCompanyEmployees(ctx, companyID)
	if err != nil {
		return models.Company{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company employees, cause: %w", err),
		)
	}

	var recipientIDs []uuid.UUID
	for _, emp := range employees.Data {
		recipientIDs = append(recipientIDs, emp.UserID)
	}

	now := time.Now().UTC()
	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		if err = s.db.UpdateCompaniesStatus(ctx, companyID, status, now); err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update company status, cause: %w", err),
			)
		}

		switch status {
		case enum.CompanyStatusActive:
			if err = s.event.PublishCompanyActivated(ctx, company, recipientIDs...); err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to publish company unblocked event, cause: %w", err),
				)
			}

		case enum.CompanyStatusInactive:
			if err = s.event.PublishCompanyDeactivated(ctx, company, recipientIDs...); err != nil {
				return errx.ErrorInternal.Raise(
					fmt.Errorf("failed to publish company blocked event, cause: %w", err),
				)
			}
		}

		return nil
	}); err != nil {
		return models.Company{}, err
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
