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

func (s Service) GetDistributorEmployee(ctx context.Context, req *empProto.GetDistributorEmployeeRequest) (*empProto.Employee, error) {
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

	distributorID, err := uuid.Parse(req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID: %s", req.DistributorId)

		return nil, errx.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid distributor ID format: %w", err),
			&errdetails.BadRequest_FieldViolation{
				Field:       "distributor_id",
				Description: "invalid distributor ID format, must be a valid UUID",
			},
		)
	}

	employee, err := s.app.GetDistributorEmployee(ctx, userID, distributorID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get employee for user ID: %s in distributor ID: %s", userID, distributorID)

		return nil, err
	}

	return responses.Employee(employee), nil
}
