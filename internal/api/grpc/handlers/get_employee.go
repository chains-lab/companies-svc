package handlers

import (
	"context"
	"fmt"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) GetEmployee(ctx context.Context, req *empProto.GetEmployeeRequest) (*empProto.Employee, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid user ID: %s", req.UserId)

		return nil, errx.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid user ID format: %w", err),
			&errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "invalid user ID format, must be a valid UUID",
			},
		)
	}

	employee, err := s.app.GetEmployee(ctx, userID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get employee for user ID: %s", userID)

		return nil, err
	}

	return responses.Employee(employee), nil
}
