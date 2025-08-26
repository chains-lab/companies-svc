package handlers

import (
	"context"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/pagi"
)

func (s Service) SelectInvites(ctx context.Context, req *empProto.SelectInvitesRequest) (*empProto.InvitesList, error) {
	filters := map[string]any{}

	if req.Filters.DistributorId != nil {
		distributorID, err := requests.DistributorID(ctx, *req.Filters.DistributorId)
		if err != nil {
			s.Log(ctx).WithError(err).Errorf("invalid distributor ID format: %s", *req.Filters.DistributorId)

			return nil, err
		}

		filters["distributor_id"] = distributorID
	}
	if req.Filters.UserId != nil {
		userID, err := requests.UserID(ctx, *req.Filters.UserId)
		if err != nil {
			s.Log(ctx).WithError(err).Errorf("invalid user ID format: %s", *req.Filters.UserId)

			return nil, err
		}

		filters["user_id"] = userID
	}
	if req.Filters.InvitedBy != nil {
		invitedByID, err := requests.InviteID(ctx, *req.Filters.InvitedBy)
		if err != nil {
			s.Log(ctx).WithError(err).Error("invalid invited by ID format")

			return nil, err
		}

		filters["invited_by"] = invitedByID
	}
	if req.Filters.Status != nil {
		status, err := requests.InviteStatus(ctx, *req.Filters.Status)
		if err != nil {
			s.Log(ctx).WithError(err).Error("invalid invite status")

			return nil, err
		}

		filters["status"] = status
	}
	if req.Filters.Role != nil {
		role, err := requests.EmployeeRole(ctx, *req.Filters.Role)
		if err != nil {
			s.Log(ctx).WithError(err).Error("invalid employee role")

			return nil, err
		}

		filters["role"] = role
	}

	ascend := true

	switch req.Sort.(type) {
	case *empProto.SelectInvitesRequest_SendAtAscend:
		ascend = true
	case *empProto.SelectInvitesRequest_SendAtDescend:
		ascend = false
	}

	invites, pag, err := s.app.SelectInvites(ctx, filters, ascend, pagi.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to select invites")

		return nil, err
	}

	return responses.EmployeeInvitesList(invites, pag), nil
}
