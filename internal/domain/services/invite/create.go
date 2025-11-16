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

func (s Service) Create(ctx context.Context, initiatorID uuid.UUID, params CreateParams) (models.Invite, error) {
	initiator, err := s.db.GetEmployeeByUserID(ctx, initiatorID)
	if err != nil {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("failed to get initiator employee by user id %s, cause: %w", initiatorID, err),
		)
	}
	if initiator.IsNil() {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployee.Raise(
			fmt.Errorf("initiator employee with user id %s not found", initiatorID),
		)
	}
	if initiator.CompanyID != params.CompanyID {
		return models.Invite{}, errx.ErrorInitiatorIsNotEmployeeInThisCompany.Raise(
			fmt.Errorf("initiator company_id %s not equal to params company_id %s", initiator.CompanyID, params.CompanyID),
		)
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

	exist, err := s.db.EmployeeExist(ctx, params.UserID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check user existence by user id %s, cause: %w", params.UserID, err),
		)
	}
	if exist {
		return models.Invite{}, errx.ErrorUserAlreadyEmployee.Raise(
			fmt.Errorf("user with id %s not found", params.UserID),
		)
	}

	company, err := s.db.GetCompanyByID(ctx, initiator.CompanyID)
	if err != nil {
		return models.Invite{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company by ID, cause: %w", err),
		)
	}
	if company.IsNil() {
		return models.Invite{}, errx.ErrorCompanyNotFound.Raise(
			fmt.Errorf("company with ID %s not found", initiator.CompanyID),
		)
	}
	if company.Status != enum.CompanyStatusActive {
		return models.Invite{}, errx.ErrorCompanyIsNotActive.Raise(
			fmt.Errorf("company with ID %s is not active", initiator.CompanyID),
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
