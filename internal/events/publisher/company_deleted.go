package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const CompanyDeletedEvent = "company.deleted"

type CompanyDeletedPayload struct {
	Company    models.Company    `json:"company"`
	Recipients PayloadRecipients `json:"recipients"`
}

func (s Service) PublishCompanyDeleted(
	ctx context.Context,
	company models.Company,
	recipients ...uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesV1,
		company.ID.String(),
		contracts.Envelope[CompanyDeletedPayload]{
			Event:     CompanyDeletedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: CompanyDeletedPayload{
				Company: company,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
