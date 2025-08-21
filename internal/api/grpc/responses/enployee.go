package responses

import (
	pagProto "github.com/chains-lab/distributors-proto/gen/go/common/pagination"
	empProto "github.com/chains-lab/distributors-proto/gen/go/svc/employee"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/pkg/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Employee(employee models.Employee) *empProto.Employee {
	return &empProto.Employee{
		UserId:        employee.UserID.String(),
		DistributorId: employee.DistributorID.String(),
		Role:          employee.Role,
		UpdatedAt:     timestamppb.New(employee.UpdatedAt),
		CreatedAt:     timestamppb.New(employee.CreatedAt),
	}
}

func EmployeesList(employees []models.Employee, pag pagination.Response) *empProto.EmployeesList {
	list := make([]*empProto.Employee, len(employees))
	for i, employee := range employees {
		list[i] = Employee(employee)
	}

	return &empProto.EmployeesList{
		Employees: list,
		Pagination: &pagProto.Response{
			Page:  pag.Page,
			Size:  pag.Size,
			Total: pag.Total,
		},
	}
}

func EmployeeInvite(invite models.Invite) *empProto.Invite {
	resp := &empProto.Invite{
		Id:            invite.ID.String(),
		DistributorId: invite.DistributorID.String(),
		UserId:        invite.UserID.String(),
		Role:          invite.Role,
		Status:        invite.Status,
		CreatedAt:     timestamppb.New(invite.CreatedAt),
	}

	if invite.AnsweredAt != nil {
		resp.AnsweredAt = timestamppb.New(*invite.AnsweredAt)
	}

	return resp
}

func EmployeeInvitesList(invites []models.Invite, pag pagination.Response) *empProto.InvitesList {
	list := make([]*empProto.Invite, len(invites))
	for i, invite := range invites {
		list[i] = EmployeeInvite(invite)
	}

	return &empProto.InvitesList{
		Invites: list,
		Pagination: &pagProto.Response{
			Page:  pag.Page,
			Size:  pag.Size,
			Total: pag.Total,
		},
	}
}
