package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) RefuseOwn(ctx context.Context, initiatorID uuid.UUID) error {
	_, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return err
	}

	err = s.employee.New().FilterUserID(initiatorID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee, cause: %w", err),
		)
	}

	return nil
}
