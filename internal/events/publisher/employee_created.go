package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const EmployeeCreatedEvent = "employee.create"

type EmployeeCreatedPayload struct {
	Company    models.Company    `json:"company"`
	Employee   models.Employee   `json:"employee"`
	Recipients PayloadRecipients `json:"recipients"`
}

func (s Service) PublishEmployeeCreated(
	ctx context.Context,
	company models.Company,
	employee models.Employee,
	recipients ...uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesV1,
		employee.CompanyID.String()+":"+employee.UserID.String(),
		contracts.Envelope[EmployeeCreatedPayload]{
			Event:     EmployeeCreatedEvent,
			Version:   "1",
			Timestamp: time.Now().UTC(),
			Data: EmployeeCreatedPayload{
				Employee: employee,
				Company:  company,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
