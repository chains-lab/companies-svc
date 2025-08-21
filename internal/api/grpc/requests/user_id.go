package requests

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/problems"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func UserID(ctx context.Context, userID string) (uuid.UUID, error) {
	res, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, problems.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid user ID format: %w", err),
			&errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "invalid user ID format, must be a valid UUID",
			},
		)
	}

	return res, nil
}
