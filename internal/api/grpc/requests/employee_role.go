package requests

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/problems"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func EmployeeRole(ctx context.Context, role string) (string, error) {
	res, err := enum.ParseEmployeeRole(role)
	if err != nil {
		return "", problems.RaiseInvalidArgument(
			ctx, err,
			&errdetails.BadRequest_FieldViolation{
				Field:       "role",
				Description: err.Error(),
			},
		)
	}

	return res, nil
}
