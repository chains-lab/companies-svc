package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorCompanyNotFound = ape.DeclareError("COMPANY_NOT_FOUND")

var ErrorCompanyIsBlocked = ape.DeclareError("COMPANY_STATUS_BLOCKED")

var ErrorCompanyIsNotActive = ape.DeclareError("COMPANY_IS_NOT_ACTIVE")

var ErrorCannotSetCompanyStatusBlocked = ape.DeclareError("COMPANY_STATUS_BLOCKED")

var ErrorCurrentEmployeeCannotCreateCompany = ape.DeclareError("CURRENT_EMPLOYEE_CANNOT_CREATE_COMPANY")

var ErrorOnlyInactiveCompanyCanBeDeleted = ape.DeclareError("COMPANY_CAN_BE_DELETED")
