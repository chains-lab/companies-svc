package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
)

const CompanyCreatedEvent = "company.created"

type CompanyCreatedPayload struct {
	Company models.Company  `json:"company"`
	Owner   models.Employee `json:"owner"`
}

func (s Service) PublishCompanyCreated(
	ctx context.Context,
	company models.Company,
	owner models.Employee,
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
				Owner:   owner,
			},
		},
	)
}
