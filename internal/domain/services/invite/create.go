package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type CreateParams struct {
	CompanyID uuid.UUID
	UserID    uuid.UUID
	Role      string
}

func (s Service) Create(ctx context.Context, initiatorUserID uuid.UUID, params CreateParams) (models.Invite, error) {
	initiator, err := s.validateInitiator(ctx, initiatorUserID, params.CompanyID)
	if err != nil {
		return models.Invite{}, err
	}

	access, err := enum.CompareEmployeeRoles(initiator.Role, params.Role)
	if err != nil {
		return models.Invite{}, errx.ErrorInvalidEmployeeRole.Raise(
			fmt.Errorf("compare employee roles: %w", err),
		)
	}
	if access <= 0 {
		return models.Invite{}, errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("initiator have not enough rights to invite role %s", params.Role),
		)
	}

	emp, err := s.db.GetEmployeeUserInCompany(ctx, params.UserID, params.CompanyID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get employee by user_id %s, cause: %w", params.UserID, err),
		)
	}
	if !emp.IsNil() {
		return models.Invite{}, errx.ErrorUserAlreadyInThisCompany.Raise(
			fmt.Errorf("employee with user_id %s already in company %s", params.UserID, params.CompanyID),
		)
	}

	company, err := s.db.GetCompanyByID(ctx, initiator.CompanyID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by EmployeeID, cause: %w", err),
		)
	}
	if company.IsNil() {
		return models.Invite{}, errx.ErrorCompanyNotFound.Raise(
			fmt.Errorf("company with EmployeeID %s not found", initiator.CompanyID),
		)
	}
	if company.Status != enum.CompanyStatusActive {
		return models.Invite{}, errx.ErrorCompanyIsNotActive.Raise(
			fmt.Errorf("company with EmployeeID %s is not active", initiator.CompanyID),
		)
	}

	inviteID := uuid.New()
	exAt := time.Now().UTC().Add(24 * time.Hour)
	now := time.Now().UTC()

	invite := models.Invite{
		ID:        inviteID,
		CompanyID: initiator.CompanyID,
		UserID:    params.UserID,
		Status:    enum.InviteStatusSent,
		Role:      params.Role,
		ExpiresAt: exAt,
		CreatedAt: now,
	}

	err = s.db.CreateInvite(ctx, invite)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create invite in db, cause: %w", err),
		)
	}

	employees, err := s.db.GetCompanyEmployees(ctx, company.ID, enum.EmployeeRoleAdmin, enum.EmployeeRoleOwner)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company employees by company id %s, cause: %w", company.ID, err),
		)
	}

	var recipientIDs []uuid.UUID
	for _, emp := range employees.Data {
		recipientIDs = append(recipientIDs, emp.UserID)
	}

	err = s.event.PublishInviteCreated(ctx, invite, company, recipientIDs...)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish invite created event, cause: %w", err),
		)
	}

	return invite, nil
}
