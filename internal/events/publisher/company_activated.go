package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const CompanyActivatedEvent = "company.activated"

type CompanyActivatedPayload struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (s Service) PublishCompanyActivated(
	ctx context.Context,
	model models.Company,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesV1,
		model.ID.String(),
		contracts.Envelope[CompanyActivatedPayload]{
			Event:     CompanyActivatedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: CompanyActivatedPayload{
				ID:        model.ID,
				Name:      model.Name,
				Icon:      model.Icon,
				Status:    model.Status,
				UpdatedAt: model.UpdatedAt,
				CreatedAt: model.CreatedAt,
			},
		},
	)
}
