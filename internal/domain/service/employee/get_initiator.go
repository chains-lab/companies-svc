package employee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/domain/errx"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetInitiator(
	ctx context.Context,
	initiatorID uuid.UUID,
) (models.Employee, error) {
	employee, err := s.employee.New().
		FilterUserID(initiatorID).
		Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Employee{}, errx.ErrorInitiatorNotEmployee.Raise(
				fmt.Errorf("initiator with userID %s not found: %w", initiatorID, err),
			)
		default:
			return models.Employee{}, errx.ErrorInternal.Raise(
				fmt.Errorf("internal error: %w", err),
			)
		}
	}

	return models.Employee{
		UserID:        employee.UserID,
		DistributorID: employee.DistributorID,
		Role:          employee.Role,
		UpdatedAt:     employee.UpdatedAt,
		CreatedAt:     employee.CreatedAt,
	}, nil
}
