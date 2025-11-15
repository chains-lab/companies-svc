package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
)

const InviteDeclinedEvent = "invite.declined"

type InviteDeclinedPayload struct {
	Invite  models.Invite  `json:"invite"`
	Company models.Company `json:"company"`
}

func (s Service) PublishInviteDeclined(
	ctx context.Context,
	invite models.Invite,
	company models.Company,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesInviteV1,
		invite.ID.String(),
		contracts.Envelope[InviteDeclinedPayload]{
			Event:     InviteDeclinedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: InviteDeclinedPayload{
				Invite:  invite,
				Company: company,
			},
		},
	)
}
