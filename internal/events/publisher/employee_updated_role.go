package publisher

import (
	"context"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/events/contracts"
	"github.com/google/uuid"
)

const EmployeeUpdatedRoleEvent = "employee.updated.role"

type EmployeeUpdatedRolePayload struct {
	Company    models.Company    `json:"company"`
	Employee   models.Employee   `json:"employee"`
	Recipients PayloadRecipients `json:"recipients"`
}

func (s Service) PublishEmployeeUpdatedRole(
	ctx context.Context,
	company models.Company,
	employee models.Employee,
	recipients []uuid.UUID,
) error {
	return s.publish(
		ctx,
		contracts.TopicCompaniesEmployeeV1,
		employee.CompanyID.String()+":"+employee.UserID.String(),
		contracts.Envelope[EmployeeUpdatedRolePayload]{
			Event:     EmployeeUpdatedRoleEvent,
			Version:   "1",
			Timestamp: employee.UpdatedAt,
			Data: EmployeeUpdatedRolePayload{
				Company:  company,
				Employee: employee,
				Recipients: PayloadRecipients{
					Users: recipients,
				},
			},
		},
	)
}
