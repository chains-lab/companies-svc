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

func (s Service) SendInvite(ctx context.Context, req *empProto.SendInviteRequest) (*empProto.Invite, error) {
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

	invite, err := s.app.SendInvite(ctx, initiator.ID, userID, distributorID, role)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to send invite for user ID: %s in distributor ID: %s", userID, distributorID)

		return nil, err
	}

	s.Log(ctx).Infof("invite sent successfully for user ID: %s in distributor ID: %s", userID, distributorID)

	return responses.EmployeeInvite(invite), nil
}
