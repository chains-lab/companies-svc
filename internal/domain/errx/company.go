package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorCompanyNotFound = ape.DeclareError("COMPANY_NOT_FOUND")

var ErrorCompanyIsBlocked = ape.DeclareError("COMPANY_STATUS_BLOCKED")

var ErrorCompanyIsNotActive = ape.DeclareError("COMPANY_IS_NOT_ACTIVE")

var ErrorOnlyInactiveCompanyCanBeDeleted = ape.DeclareError("ONLY_INACTIVE_COMPANY_CAN_BE_DELETED")

var ErrorCannotSetCompanyStatusBlocked = ape.DeclareError("CANNOT_SET_COMPANY_STATUS_BLOCKED")

var ErrorInvalidCompanyStatus = ape.DeclareError("INVALID_COMPANY_STATUS")
