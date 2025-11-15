package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const InviteAccepted = "invite.accepted"

type InviteAcceptedPayload struct {
	Invite     models.Invite     `json:"invite"`
	Company    models.Company    `json:"company"`
	Recipients PayloadRecipients `json:"recipients"`
}

func (s Service) PublishInviteAccepted(
	ctx context.Context,
	invite models.Invite,
	company models.Company,
	recipients ...uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesInviteV1,
		invite.ID.String(),
		contracts.Envelope[InviteAcceptedPayload]{
			Event:     InviteAccepted,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: InviteAcceptedPayload{
				Invite:  invite,
				Company: company,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
