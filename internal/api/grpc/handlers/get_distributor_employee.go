package handlers

import (
	"context"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) GetDistributorEmployee(ctx context.Context, req *empProto.GetDistributorEmployeeRequest) (*empProto.Employee, error) {
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

	employee, err := s.app.GetDistributorEmployee(ctx, userID, distributorID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get employee for user ID: %s in distributor ID: %s", userID, distributorID)

		return nil, err
	}

	return responses.Employee(employee), nil
}
