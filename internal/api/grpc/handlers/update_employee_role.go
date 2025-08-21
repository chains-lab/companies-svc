package handlers

import (
	"context"
	"fmt"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) UpdateEmployeeRole(ctx context.Context, req *empProto.UpdateEmployeeRoleRequest) (*empProto.Employee, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid user ID: %s", req.UserId)
		return nil, errx.RaiseInvalidArgument(
			ctx, err,
			&errdetails.BadRequest_FieldViolation{
				Field:       "user_id",
				Description: "invalid UUID format for user ID",
			},
		)
	}

	distributorID, err := uuid.Parse(req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID: %s", req.DistributorId)

		return nil, errx.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid distributor ID: %w", err),

			&errdetails.BadRequest_FieldViolation{
				Field:       "distributor_id",
				Description: "invalid UUID format for distributor ID",
			},
		)
	}

	role, err := enum.ParseEmployeeRole(req.Role)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid employee role: %s", req.Role)

		return nil, errx.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid employee role: %w", err),
			&errdetails.BadRequest_FieldViolation{
				Field:       "role",
				Description: "invalid employee role",
			},
		)
	}

	res, err := s.app.UpdateEmployeeRole(ctx, initiator.ID, userID, distributorID, role)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to update employee role for ID: %s", userID)
		return nil, err
	}

	s.Log(ctx).Infof("employee %s role updated successfully to %s", userID, req.Role)

	return responses.Employee(res), nil
}
