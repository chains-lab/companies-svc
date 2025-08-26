package handlers

import (
	"context"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s Service) DeleteEmployee(ctx context.Context, req *empProto.DeleteEmployeeRequest) (*emptypb.Empty, error) {
	initiator, err := meta.User(ctx)
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to get user from context")
		return nil, err
	}

	userID, err := requests.UserID(ctx, req.UserId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return nil, err
	}

	distributorID, err := requests.DistributorID(ctx, req.DistributorId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid distributor ID format: %s", req.DistributorId)

		return nil, err
	}

	err = s.app.DeleteEmployee(ctx, initiator.ID, userID, distributorID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to delete employee with user_id: %s", userID)
		return nil, err
	}

	s.Log(ctx).Infof("employee %s deleted successfully", userID)

	return &emptypb.Empty{}, nil
}
