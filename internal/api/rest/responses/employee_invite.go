package responses

import (
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/resources"
)

func EmployeeInvites(m models.Invite) resources.Invite {
	resp := resources.Invite{
		Data: resources.InviteData{
			Id:   m.ID.String(),
			Type: resources.InviteType,
			Attributes: resources.InviteAttributes{
				Status:        m.Status,
				Role:          m.Role,
				DistributorId: m.DistributorID.String(),
				Token:         m.Token,
				ExpiresAt:     m.ExpiresAt,
				CreatedAt:     m.CreatedAt,
			},
		},
	}

	return resp
}
