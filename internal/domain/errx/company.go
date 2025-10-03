package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorcompanyNotFound = ape.DeclareError("company_NOT_FOUND")

var ErrorcompanyIsBlocked = ape.DeclareError("company_STATUS_BLOCKED")

var ErrorcompanyIsNotActive = ape.DeclareError("company_IS_NOT_ACTIVE")

var ErrorCannotSetcompaniestatusBlocked = ape.DeclareError("CANNOT_SET_company_STATUS_BLOCKED")

var ErrorCurrentEmployeeCannotCreatecompany = ape.DeclareError("CURRENT_EMPLOYEE_CANNOT_CREATE_company")
