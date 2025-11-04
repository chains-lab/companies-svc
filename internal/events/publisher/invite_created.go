package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const InviteCreated = "invite.created"

type InviteCreatedPayload struct {
	ID        uuid.UUID `json:"id"`
	CompanyID uuid.UUID `json:"company_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (s Service) PublishInviteCreated(
	ctx context.Context,
	invite models.Invite,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesInviteV1,
		invite.ID.String(),
		contracts.Envelope[InviteCreatedPayload]{
			Event:     InviteCreated,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: InviteCreatedPayload{
				ID:        invite.ID,
				CompanyID: invite.CompanyID,
				UserID:    invite.UserID,
				Role:      invite.Role,
				Status:    invite.Status,
				ExpiresAt: invite.ExpiresAt,
				CreatedAt: invite.CreatedAt,
			},
		},
	)
}
