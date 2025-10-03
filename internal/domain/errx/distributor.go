package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorDistributorNotFound = ape.DeclareError("DISTRIBUTOR_NOT_FOUND")

var ErrorDistributorIsBlocked = ape.DeclareError("DISTRIBUTOR_STATUS_BLOCKED")

var ErrorDistributorIsNotActive = ape.DeclareError("DISTRIBUTOR_IS_NOT_ACTIVE")

var ErrorCannotSetDistributorStatusBlocked = ape.DeclareError("CANNOT_SET_DISTRIBUTOR_STATUS_BLOCKED")

var ErrorCurrentEmployeeCannotCreateDistributor = ape.DeclareError("CURRENT_EMPLOYEE_CANNOT_CREATE_DISTRIBUTOR")
