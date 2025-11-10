package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const CompanyUnblockedEvent = "company.unblocked"

type CompanyUnblockedPayload struct {
	Block      models.CompanyBlock `json:"block"`
	Company    models.Company      `json:"company"`
	Recipients PayloadRecipients   `json:"recipients"`
}

func (s Service) PublishCompanyUnblocked(
	ctx context.Context,
	block models.CompanyBlock,
	company models.Company,
	recipients []uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesV1,
		company.ID.String(),
		contracts.Envelope[CompanyUnblockedPayload]{
			Event:     CompanyUnblockedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: CompanyUnblockedPayload{
				Company: company,
				Block:   block,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
