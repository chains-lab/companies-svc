package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) DeleteByInitiatorID(
	ctx context.Context,
	initiatorID, userID, companyID uuid.UUID,
) error {
	employee, err := s.Get(ctx, GetParams{
		UserID:    &userID,
		CompanyID: &companyID,
	})
	if err != nil {
		return err
	}

	initiator, err := s.validateInitiatorRight(ctx, initiatorID, &companyID, enum.EmployeeRoleOwner, enum.EmployeeRoleAdmin)
	if err != nil {
		return err
	}
	if initiator.UserID == userID {
		return errx.ErrorCannotDeleteYourself.Raise(
			fmt.Errorf("initiator %s is trying to delete himself", initiatorID),
		)
	}

	company, err := s.getCompany(ctx, companyID)
	if err != nil {
		return err
	}
	if company.Status != enum.CompanyStatusActive {
		return errx.ErrorCompanyIsNotActive.Raise(
			fmt.Errorf("cannot delete employee from inactive company with ID %s", companyID),
		)
	}

	allowed, err := enum.CompareEmployeeRoles(initiator.Role, employee.Role)
	if err != nil {
		return errx.ErrorInvalidEmployeeRole.Raise(err)
	}
	if allowed != 1 {
		return errx.ErrorInitiatorHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator have not enough rights to delete employee"),
		)
	}

	return s.delete(ctx, employee, company)
}

func (s Service) DeleteMe(
	ctx context.Context,
	initiatorID uuid.UUID,
) error {
	own, err := s.GetInitiator(ctx, initiatorID)
	if err != nil {
		return err
	}

	if own.Role == enum.EmployeeRoleOwner {
		return errx.ErrorOwnerCannotRefuseSelf.Raise(
			fmt.Errorf("owner cannot refuse self"),
		)
	}

	company, err := s.getCompany(ctx, own.CompanyID)
	if err != nil {
		return err
	}
	//if company.Status != enum.CompanyStatusActive {
	//	return errx.ErrorCompanyIsNotActive.Raise(
	//		fmt.Errorf("cannot refuse employee from inactive company with ID %s", own.CompanyID),
	//	)
	//}

	return s.delete(ctx, own, company)
}

func (s Service) delete(
	ctx context.Context,
	employee models.Employee,
	company models.Company,
) error {
	employees, err := s.db.GetCompanyEmployees(ctx, company.ID, enum.EmployeeRoleAdmin, enum.EmployeeRoleOwner)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company employees, cause: %w", err),
		)
	}

	var recipientsIDs []uuid.UUID
	for _, emp := range employees.Data {
		recipientsIDs = append(recipientsIDs, emp.UserID)
	}

	if err = s.db.DeleteEmployee(ctx, employee.UserID, company.ID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete employee, cause: %w", err),
		)
	}

	if err = s.event.PublishEmployeeDeleted(ctx, company, employee, recipientsIDs...); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish company deleted event, cause: %w", err),
		)
	}

	return nil
}
