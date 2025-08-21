package handlers

import (
	"context"
	"fmt"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/distributors-svc/pkg/pagination"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (s Service) SelectInvites(ctx context.Context, req *empProto.SelectInvitesRequest) (*empProto.InvitesList, error) {
	filters := map[string]any{}

	if req.Filters.DistributorId != nil {
		distributorID, err := uuid.Parse(*req.Filters.DistributorId)
		if err != nil {
			s.Log(ctx).WithError(err).Error("invalid distributor ID format")
			return nil, errx.RaiseInvalidArgument(
				ctx, fmt.Errorf("invalid distributor ID format: %w", err),
				&errdetails.BadRequest_FieldViolation{
					Field:       "filters.distributor_id",
					Description: "invalid UUID format for distributor ID",
				},
			)
		}

		filters["distributor_id"] = distributorID
	}
	if req.Filters.UserId != nil {
		userID, err := uuid.Parse(*req.Filters.UserId)
		if err != nil {
			s.Log(ctx).WithError(err).Error("invalid user ID format")
			return nil, errx.RaiseInvalidArgument(
				ctx, fmt.Errorf("invalid user ID format: %w", err),
				&errdetails.BadRequest_FieldViolation{
					Field:       "filters.user_id",
					Description: "invalid UUID format for user ID",
				},
			)
		}

		filters["user_id"] = userID
	}
	if req.Filters.InvitedBy != nil {
		invitedByID, err := uuid.Parse(*req.Filters.InvitedBy)
		if err != nil {
			s.Log(ctx).WithError(err).Error("invalid invited by ID format")
			return nil, errx.RaiseInvalidArgument(
				ctx, fmt.Errorf("invalid invited by ID format: %w", err),
				&errdetails.BadRequest_FieldViolation{
					Field:       "filters.invited_by",
					Description: "invalid UUID format for invited by ID",
				},
			)
		}

		filters["invited_by"] = invitedByID
	}
	if req.Filters.Status != nil {
		status, err := enum.ParseInviteStatus(*req.Filters.Status)
		if err != nil {
			s.Log(ctx).WithError(err).Error("invalid invite status")
			return nil, errx.RaiseInvalidArgument(
				ctx, fmt.Errorf("invalid invite status: %w", err),
				&errdetails.BadRequest_FieldViolation{
					Field:       "filters.status",
					Description: "invalid invite status",
				},
			)
		}

		filters["status"] = status
	}
	if req.Filters.Role != nil {
		role, err := enum.ParseEmployeeRole(*req.Filters.Role)
		if err != nil {
			s.Log(ctx).WithError(err).Error("invalid employee role")
			return nil, errx.RaiseInvalidArgument(
				ctx, fmt.Errorf("invalid employee role: %w", err),
				&errdetails.BadRequest_FieldViolation{
					Field:       "filters.role",
					Description: "invalid employee role",
				},
			)
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

	invites, pag, err := s.app.SelectInvites(ctx, filters, ascend, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})

	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to select invites")

		return nil, errx.RaiseInternal(ctx, fmt.Errorf("selecting invites: %w", err))
	}

	return responses.EmployeeInvitesList(invites, pag), nil
}
