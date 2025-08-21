package handlers

import (
	"context"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) GetInvite(ctx context.Context, req *empProto.GetInviteRequest) (*empProto.Invite, error) {
	inviteID, err := requests.InviteID(ctx, req.InviteId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid invite ID format: %s", req.InviteId)

		return nil, err
	}

	invite, err := s.app.GetInvite(ctx, inviteID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get invite with ID: %s", inviteID)

		return nil, err
	}

	return responses.EmployeeInvite(invite), nil
}
