package handlers

import (
	"context"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) SendInvite(ctx context.Context, req *empProto.SendInviteRequest) (*empProto.Invite, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")

		return nil, err
	}

	userID, err := requests.UserID(ctx, req.UserId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid user ID: %s", req.UserId)

		return nil, err
	}

	distributorID, err := requests.DistributorID(ctx, req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID: %s", req.DistributorId)

		return nil, err
	}

	role, err := requests.EmployeeRole(ctx, req.Role)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid employee role: %s", req.Role)

		return nil, err
	}

	invite, err := s.app.SendInvite(ctx, initiator.ID, userID, distributorID, role)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to send invite for user ID: %s in distributor ID: %s", userID, distributorID)

		return nil, err
	}

	s.Log(ctx).Infof("invite sent successfully for user ID: %s in distributor ID: %s", userID, distributorID)

	return responses.EmployeeInvite(invite), nil
}
