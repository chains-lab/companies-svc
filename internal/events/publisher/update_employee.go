package publisher

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/events"
	"github.com/google/uuid"
)

type EmployeeUpdate struct {
	UserID    uuid.UUID  `json:"user_id"`
	CompanyID *uuid.UUID `json:"company_id"`
	Role      *string    `json:"role"`
}

func (s *Service) UpdateEmployee(
	ctx context.Context,
	userID uuid.UUID,
	companyID *uuid.UUID,
	role *string,
) error {
	env := events.Envelope[EmployeeUpdate]{
		Event:     "employee.update",
		Version:   "1",
		Timestamp: time.Now().UTC(),
		Data: EmployeeUpdate{
			UserID:    userID,
			CompanyID: companyID,
			Role:      role,
		},
	}
	return s.publish(
		ctx,
		events.TopicCompaniesEmployeeV1,
		userID.String(),
		env,
	)
}
