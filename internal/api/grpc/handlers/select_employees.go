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

func (s Service) SelectEmployees(ctx context.Context, req *empProto.SelectEmployeesRequest) (*empProto.EmployeesList, error) {
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
	case *empProto.SelectEmployeesRequest_RolesAscend:
		ascend = true
	case *empProto.SelectEmployeesRequest_RolesDescend:
		ascend = false
	}

	employees, pag, err := s.app.SelectEmployees(ctx, filters, ascend, pagination.Request{
		Page: req.Pagination.Page,
		Size: req.Pagination.Size,
	})
	if err != nil {
		s.Log(ctx).WithError(err).Error("failed to select employees")

		return nil, errx.RaiseInternal(ctx, fmt.Errorf("selecting employees: %w", err))
	}

	return responses.EmployeesList(employees, pag), nil
}
