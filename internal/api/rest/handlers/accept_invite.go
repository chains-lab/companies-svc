package handlers

import (
	"net/http"

	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
)

func (s Service) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")
		return nil, err
	}

	inviteID, err := requests.InviteID(ctx, req.InviteId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid invite ID format: %s", req.InviteId)

		return nil, err
	}

	invite, err := s.app.AcceptInvite(ctx, initiator.ID, inviteID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to accept invite with ID: %s", req.InviteId)
		return nil, err
	}

	return responses.EmployeeInvite(invite), nil
}
