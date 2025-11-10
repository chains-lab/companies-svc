package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const CompanyBlockedEvent = "company.blocked"

type CompanyBlockedPayload struct {
	Block      models.CompanyBlock `json:"block"`
	Company    models.Company      `json:"company"`
	Recipients PayloadRecipients   `json:"recipients"`
}

func (s Service) PublishCompanyBlocked(
	ctx context.Context,
	block models.CompanyBlock,
	company models.Company,
	recipients []uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesV1,
		company.ID.String(),
		contracts.Envelope[CompanyBlockedPayload]{
			Event:     CompanyBlockedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: CompanyBlockedPayload{
				Block:   block,
				Company: company,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
