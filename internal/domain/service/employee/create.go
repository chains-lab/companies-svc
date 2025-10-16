package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

type CreateParams struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Role      string
}

func (s Service) Create(ctx context.Context, params CreateParams) (models.EmployeeWithUserData, error) {
	emp, err := s.db.GetEmployeeByUserID(ctx, params.UserID)
	if err != nil {
		return models.EmployeeWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("falied to check existing employee, cause: %w", err),
		)
	}

	if !emp.IsNil() {
		return models.EmployeeWithUserData{}, errx.ErrorEmployeeAlreadyExists.Raise(
			fmt.Errorf("employee with user ID %s already exists", params.UserID),
		)
	}

	now := time.Now().UTC()
	err = enum.CheckEmployeeRole(params.Role)
	if err != nil {
		return models.EmployeeWithUserData{}, errx.ErrorInvalidEmployeeRole.Raise(
			fmt.Errorf("failed to check employee role, cause: %w", err),
		)
	}

	emp = models.Employee{
		UserID:    params.UserID,
		CompanyID: params.CompanyID,
		Role:      params.Role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.db.CreateEmployee(ctx, emp)
	if err != nil {
		return models.EmployeeWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create employee, cause: %w", err),
		)
	}

	profiles, err := s.userGuesser.Guess(ctx, params.UserID)
	if err != nil {
		return models.EmployeeWithUserData{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to guess employee, cause: %w", err),
		)
	}

	return emp.AddProfileData(profiles[params.UserID]), nil
}
