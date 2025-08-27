package responses

import (
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/resources"
	"github.com/chains-lab/pagi"
)

func EmployeeInvites(m models.EmployeeInvite) resources.EmployeeInvite {
	resp := resources.EmployeeInvite{
		Data: resources.EmployeeInviteData{
			Id:   m.ID.String(),
			Type: resources.EmployeeInviteType,
			Attributes: resources.EmployeeInviteAttributes{
				UserId:        m.UserID.String(),
				DistributorId: m.DistributorID.String(),
				InvitedBy:     m.InvitedBy.String(),
				Role:          m.Role,
				Status:        m.Status,
				CreatedAt:     m.CreatedAt,
			},
		},
	}

	if m.AnsweredAt != nil {
		resp.Data.Attributes.AnsweredAt = m.AnsweredAt
	}

	return resp
}

func EmployeeInvitesCollection(ms []models.EmployeeInvite, pag pagi.Response) resources.EmployeesInvitesCollection {
	items := make([]resources.EmployeeInviteData, 0, len(ms))
	for _, m := range ms {
		items = append(items, EmployeeInvites(m).Data)
	}

	return resources.EmployeesInvitesCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}
}
