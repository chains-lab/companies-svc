package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) RefuseMe(ctx context.Context, initiatorID uuid.UUID) error {
	own, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return err
	}

	if own.Role == enum.EmployeeRoleOwner {
		return errx.ErrorOwnerCannotRefuseSelf.Raise(
			fmt.Errorf("owner cannot refuse self"),
		)
	}
	err = s.db.DeleteEmployee(ctx, own.UserID, own.CompanyID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee, cause: %w", err),
		)
	}

	return nil
}
