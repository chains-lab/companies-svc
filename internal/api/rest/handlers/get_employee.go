package handlers

import (
	"context"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
)

func (s Service) GetEmployee(ctx context.Context, req *empProto.GetEmployeeRequest) (*empProto.Employee, error) {
	userID, err := requests.UserID(ctx, req.UserId)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("invalid user ID format: %s", req.UserId)

		return nil, err
	}

	employee, err := s.app.GetEmployee(ctx, userID)
	if err != nil {
		s.Log(ctx).WithError(err).Errorf("failed to get employee for user ID: %s", userID)

		return nil, err
	}

	return responses.Employee(employee), nil
}
