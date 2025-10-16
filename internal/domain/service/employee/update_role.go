package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) UpdateEmployeeRole(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	newRole string,
) (models.EmployeeWithUserData, error) {
	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return models.EmployeeWithUserData{}, err
	}

	user, err := s.Get(ctx, GetParams{
		UserID: &userID,
	})
	if err != nil {
		return models.EmployeeWithUserData{}, err
	}

	if initiator.CompanyID != user.CompanyID {
		return models.EmployeeWithUserData{}, errx.ErrorInitiatorIsNotEmployeeOfThiscompany.Raise(
			fmt.Errorf("initiator %s and chosen employee %s have different companies", initiatorID, userID),
		)
	}

	allowed, err := enum.CompareEmployeeRoles(initiator.Role, user.Role)
	if err != nil {
		return models.EmployeeWithUserData{}, errx.EmployeeInvalidRole.Raise(
			fmt.Errorf("new role is invalid: %w", err),
		)
	}
	if allowed != 1 {
		return models.EmployeeWithUserData{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	allowed, err = enum.CompareEmployeeRoles(initiator.Role, newRole)
	if err != nil {
		return models.EmployeeWithUserData{}, errx.EmployeeInvalidRole.Raise(
			fmt.Errorf("new role is invalid: %w", err),
		)
	}
	if allowed != 1 {
		return models.EmployeeWithUserData{}, errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to update employee role"),
		)
	}

	now := time.Now().UTC()

	err = s.db.UpdateEmployeeRole(ctx, userID, newRole, now)
	if err != nil {
		return models.EmployeeWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update employee role, cause: %w", err),
		)
	}

	profiles, err := s.userGuesser.Guess(ctx, userID)
	if err != nil {
		return models.EmployeeWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to guess employee, cause: %w", err),
		)
	}

	return models.EmployeeWithUserData{
		UserID:    user.UserID,
		Username:  profiles[user.UserID].Username,
		Avatar:    profiles[user.UserID].Avatar,
		CompanyID: user.CompanyID,
		Role:      newRole,
		UpdatedAt: now,
		CreatedAt: user.CreatedAt,
	}, nil
}
