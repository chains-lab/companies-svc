package handlers

import (
	"context"
	"fmt"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteEmployee(ctx context.Context, req *empProto.DeleteEmployeeRequest) (*emptypb.Empty, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")
		return nil, err
	}

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

	err = s.app.DeleteEmployee(ctx, initiator.ID, userID, distributorID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to delete employee with user_id: %s", userID)
		return nil, err
	}

	s.Log(ctx).Infof("employee %s deleted successfully", userID)

	return &emptypb.Empty{}, nil
}
