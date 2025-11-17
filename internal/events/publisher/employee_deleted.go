package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const EmployeeDeletedEvent = "employee.delete"

type EmployeeDeletedPayload struct {
	Company    models.Company    `json:"company"`
	Employee   models.Employee   `json:"employee"`
	Recipients PayloadRecipients `json:"recipients"`
}

func (s Service) PublishEmployeeDeleted(
	ctx context.Context,
	company models.Company,
	employee models.Employee,
	recipients ...uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesEmployeeV1,
		employee.ID.String(),
		contracts.Envelope[EmployeeDeletedPayload]{
			Event:     EmployeeDeletedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: EmployeeDeletedPayload{
				Company:  company,
				Employee: employee,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
