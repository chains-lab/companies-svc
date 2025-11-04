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
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	Role      string    `json:"role"`
	Position  *string   `json:"position"`
	Label     *string   `json:"label,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (s Service) PublishEmployeeUpdated(
	ctx context.Context,
	employee models.Employee,
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
				UserID:    employee.UserID,
				CompanyID: employee.CompanyID,
				Role:      employee.Role,
				Position:  employee.Position,
				Label:     employee.Label,
				UpdatedAt: employee.UpdatedAt,
				CreatedAt: employee.CreatedAt,
			},
		},
	)
}
