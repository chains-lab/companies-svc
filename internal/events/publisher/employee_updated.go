package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const EmployeeUpdatedEvent = "employee.update"

type EmployeeUpdatedPayload struct {
	Company    models.Company    `json:"company"`
	Employee   models.Employee   `json:"employee"`
	Recipients PayloadRecipients `json:"recipients"`
}

func (s Service) PublishEmployeeUpdated(
	ctx context.Context,
	company models.Company,
	employee models.Employee,
	recipients ...uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesEmployeeV1,
		employee.CompanyID.String()+":"+employee.UserID.String(),
		contracts.Envelope[EmployeeUpdatedPayload]{
			Event:     EmployeeUpdatedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: EmployeeUpdatedPayload{
				Company:  company,
				Employee: employee,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
