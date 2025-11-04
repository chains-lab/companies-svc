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
	ID          uuid.UUID  `json:"id"`
	CompanyID   uuid.UUID  `json:"company_id"`
	InitiatorID uuid.UUID  `json:"initiator_id"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	BlockedAt   time.Time  `json:"blocked_at"`
	CanceledAt  *time.Time `json:"canceled_at"`
}

func (s Service) PublishCompanyUnblocked(
	ctx context.Context,
	block models.CompanyBlock,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesCompanyBlockV1,
		block.ID.String(),
		contracts.Envelope[CompanyUnblockedPayload]{
			Event:     CompanyUnblockedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: CompanyUnblockedPayload{
				ID:          block.ID,
				CompanyID:   block.CompanyID,
				InitiatorID: block.InitiatorID,
				Reason:      block.Reason,
				Status:      block.Status,
				BlockedAt:   block.BlockedAt,
				CanceledAt:  block.CanceledAt,
			},
		},
	)
}
