package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const CompanyCreatedEvent = "company.created"

type CompanyCreatedPayload struct {
	Company models.Company `json:"company"`
	OwnerID uuid.UUID      `json:"owner_id"`
}

func (s Service) PublishCompanyCreated(
	ctx context.Context,
	company models.Company,
	ownerID uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesV1,
		company.ID.String(),
		contracts.Envelope[CompanyCreatedPayload]{
			Event:     CompanyCreatedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: CompanyCreatedPayload{
				Company: company,
				OwnerID: ownerID,
			},
		},
	)
}
