package handlers

import (
	"context"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) UpdateEmployeeRole(ctx context.Context, req *empProto.UpdateEmployeeRoleRequest) (*empProto.Employee, error) {
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

	res, err := s.app.UpdateEmployeeRole(ctx, initiator.ID, userID, distributorID, role)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to update employee role for ID: %s", userID)
		return nil, err
	}

	s.Log(ctx).Infof("employee %s role updated successfully to %s", userID, req.Role)

	return responses.Employee(res), nil
}
