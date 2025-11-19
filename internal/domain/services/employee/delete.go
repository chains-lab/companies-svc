package employee

import (
	"context"
	"fmt"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/errx"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) DeleteByEmployee(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	companyID uuid.UUID,
) error {
	employee, err := s.Get(ctx, userID, companyID)
	if err != nil {
		return err
	}

	if employee.UserID == initiatorID {
		return errx.ErrorCannotDeleteYourself.Raise(
			fmt.Errorf("initiator %s is trying to delete himself", initiatorID),
		)
	}

	initiator, err := s.validateInitiator(
		ctx, initiatorID, companyID,
		enum.EmployeeRoleOwner, enum.EmployeeRoleAdmin,
	)
	if err != nil {
		return err
	}

	company, err := s.getCompany(ctx, employee.CompanyID)
	if err != nil {
		return err
	}
	if company.Status != enum.CompanyStatusActive {
		return errx.ErrorCompanyIsNotActive.Raise(
			fmt.Errorf("cannot delete employee from inactive company with EmployeeID %s", employee.CompanyID),
		)
	}

	allowed, err := enum.CompareEmployeeRoles(initiator.Role, employee.Role)
	if err != nil {
		return errx.ErrorInvalidEmployeeRole.Raise(err)
	}
	if allowed != 1 {
		return errx.ErrorNotEnoughRight.Raise(
			fmt.Errorf("initiator have not enough rights to delete employee"),
		)
	}

	return s.delete(ctx, employee, company)
}

func (s Service) DeleteMe(
	ctx context.Context,
	userID uuid.UUID,
	companyID uuid.UUID,
) error {
	own, err := s.Get(ctx, userID, companyID)
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
	//		fmt.Errorf("cannot refuse employee from inactive company with EmployeeID %s", own.CompanyID),
	//	)
	//}

	return s.delete(ctx, own, company)
}

func (s Service) delete(
	ctx context.Context,
	employee models.Employee,
	company models.Company,
) error {
	if err := s.db.DeleteEmployee(ctx, employee.UserID, employee.CompanyID); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete employee, cause: %w", err),
		)
	}

	employees, err := s.db.GetCompanyEmployees(ctx, company.ID, enum.EmployeeRoleAdmin, enum.EmployeeRoleOwner)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get company employees, cause: %w", err),
		)
	}

	if err = s.event.PublishEmployeeDeleted(ctx, company, employee, employees.GetUserIDs()...); err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish company deleted event, cause: %w", err),
		)
	}

	return nil
}
