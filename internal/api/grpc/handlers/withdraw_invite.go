package handlers

import (
	"context"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) WithdrewInvite(ctx context.Context, req *empProto.WithdrawInviteRequest) (*empProto.Invite, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")
		return nil, err
	}

	inviteID, err := uuid.Parse(req.InviteId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid invite ID: %s", req.InviteId)
		return nil, errx.RaiseInvalidArgument(
			ctx, err,
			&errdetails.BadRequest_FieldViolation{
				Field:       "id",
				Description: "invalid UUID format for invite ID",
			},
		)
	}

	invite, err := s.app.WithdrawInvite(ctx, initiator.ID, inviteID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to withdraw invite with ID: %s", req.InviteId)

		return nil, err
	}

	return responses.EmployeeInvite(invite), nil
}
