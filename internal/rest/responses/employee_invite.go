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
				Status: m.Status,
				Role:   m.Role,
				//UserId:    m.UserID,
				Token:     m.Token,
				ExpiresAt: m.ExpiresAt,
				CreatedAt: m.CreatedAt,
			},
		},
	}

	return resp
}
