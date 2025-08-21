package requests

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/problems"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func InviteStatus(ctx context.Context, status string) (string, error) {
	res, err := enum.ParseInviteStatus(status)
	if err != nil {
		return "", problems.RaiseInvalidArgument(
			ctx, err,
			&errdetails.BadRequest_FieldViolation{
				Field:       "status",
				Description: err.Error(),
			},
		)
	}

	return res, nil
}
