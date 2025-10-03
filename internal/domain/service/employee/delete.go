package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, initiatorID, userID, companyID uuid.UUID) error {
	user, err := s.Get(ctx, GetParams{
		UserID:    &userID,
		CompanyID: &companyID,
	})
	if err != nil {
		return err
	}

	initiator, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return err
	}

	if initiator.UserID == userID {
		return errx.ErrorCannotDeleteYourself.Raise(
			fmt.Errorf("initiator %s is trying to delete himself", initiatorID),
		)
	}

	if initiator.CompanyID != user.CompanyID {
		return errx.ErrorInitiatorIsNotEmployeeOfThiscompany.Raise(
			fmt.Errorf("initiator %s and chosen employee %s have different companies", initiatorID, userID),
		)
	}

	allowed, err := enum.CompareEmployeeRoles(initiator.Role, user.Role)
	if err != nil {
		return errx.ErrorInvalidEmployeeRole.Raise(err)
	}
	if allowed != 1 {
		return errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to delete employee"),
		)
	}

	err = s.db.DeleteEmployee(ctx, userID, companyID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete employee, cause: %w", err),
		)
	}

	return nil
}
