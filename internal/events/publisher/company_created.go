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
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (s Service) PublishCompanyCreated(
	ctx context.Context,
	company models.Company,
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
				ID:        company.ID,
				Name:      company.Name,
				Icon:      company.Icon,
				Status:    company.Status,
				UpdatedAt: company.UpdatedAt,
				CreatedAt: company.CreatedAt,
			},
		},
	)
}
