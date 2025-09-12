package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
)

func (e Employee) RefuseOwn(ctx context.Context, initiatorID uuid.UUID) error {
	_, err := e.GetInitiator(ctx, initiatorID)
	if err != nil {
		return err
	}

	err = e.employee.New().FilterUserID(initiatorID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to refuse own employee, cause: %w", err),
		)
	}

	return nil
}
