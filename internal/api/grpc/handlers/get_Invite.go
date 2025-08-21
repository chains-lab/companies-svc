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

func (s Service) GetInvite(ctx context.Context, req *empProto.GetInviteRequest) (*empProto.Invite, error) {
	inviteID, err := uuid.Parse(req.InviteId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid invite ID: %s", req.InviteId)

		return nil, errx.RaiseInvalidArgument(
			ctx, fmt.Errorf("invalid invite ID format: %w", err),
			&errdetails.BadRequest_FieldViolation{
				Field:       "invite_id",
				Description: "invalid invite ID format, must be a valid UUID",
			},
		)
	}

	invite, err := s.app.GetInvite(ctx, inviteID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get invite with ID: %s", inviteID)

		return nil, err
	}

	return responses.EmployeeInvite(invite), nil
}
