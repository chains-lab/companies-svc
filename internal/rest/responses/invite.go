package responses

import (
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/resources"
)

func Invites(m models.Invite) resources.Invite {
	resp := resources.Invite{
		Data: resources.InviteData{
			Id:   m.ID,
			Type: resources.InviteType,
			Attributes: resources.InviteAttributes{
				UserId:    m.UserID,
				CompanyId: m.CompanyID,
				Status:    m.Status,
				Role:      m.Role,
				ExpiresAt: m.ExpiresAt,
				CreatedAt: m.CreatedAt,
			},
		},
	}

	return resp
}
