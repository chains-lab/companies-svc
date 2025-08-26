package handlers

import (
	"context"
	"fmt"

	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/requests"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/responses"
	"github.com/chains-lab/distributors-svc/internal/problems"
	"github.com/chains-lab/pagi"
)

func (s Service) SelectEmployees(ctx context.Context, req *empProto.SelectEmployeesRequest) (*empProto.EmployeesList, error) {
	filters := map[string]any{}

	if req.Filters.DistributorId != nil {
		distributorID, err := requests.DistributorID(ctx, *req.Filters.DistributorId)
		if err != nil {
			s.Log(ctx).WithError(err).Errorf("invalid distributor ID format: %s", *req.Filters.DistributorId)

			return nil, err
		}

		filters["distributor_id"] = distributorID
	}
	if req.Filters.Role != nil {
		role, err := requests.EmployeeRole(ctx, *req.Filters.Role)
		if err != nil {
			s.Log(ctx).WithError(err).Errorf("invalid employee role: %s", *req.Filters.Role)

			return nil, err
		}

		filters["role"] = role
	}

	ascend := true

	switch req.Sort.(type) {
	case *empProto.SelectEmployeesRequest_RolesAscend:
		ascend = true
	case *empProto.SelectEmployeesRequest_RolesDescend:
		ascend = false
	}

	employees, pag, err := s.app.SelectEmployees(ctx, filters, ascend, pagi.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to select employees")

		return nil, problems.RaiseInternal(ctx, fmt.Errorf("selecting employees: %w", err))
	}

	return responses.EmployeesList(employees, pag), nil
}
