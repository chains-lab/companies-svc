package requests

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/problems"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func DistributorID(ctx context.Context, distributorID string) (uuid.UUID, error) {
	res, err := uuid.Parse(distributorID)
	if err != nil {
		return uuid.Nil, problems.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid distributor ID format: %w", err),
			&errdetails.BadRequest_FieldViolation{
				Field:       "distributor_id",
				Description: "invalid distributor ID format, must be a valid UUID",
			},
		)
	}

	return res, nil
}
