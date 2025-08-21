package requests

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/problems"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func BlockID(ctx context.Context, blockID string) (uuid.UUID, error) {
	res, err := uuid.Parse(blockID)
	if err != nil {
		return uuid.Nil, problems.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid block ID format: %w", err),
			&errdetails.BadRequest_FieldViolation{
				Field:       "block_id",
				Description: "invalid block ID format, must be a valid UUID",
			},
		)
	}

	return res, nil
}
