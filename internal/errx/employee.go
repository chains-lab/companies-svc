package errx

import (
	"github.com/chains-lab/ape"
)

var EmployeeNotFound = ape.DeclareError("EMPLOYEE_NOT_FOUND")

var EmployeeAlreadyExists = ape.DeclareError("EMPLOYEE_ALREADY_EXISTS")

var EmployeeRoleNotSupported = ape.DeclareError("EMPLOYEE_ROLE_NOT_SUPPORTED")

var InitiatorNotEmployee = ape.DeclareError("INITIATOR_NOT_EMPLOYEE")

var InitiatorIsAlreadyEmployee = ape.DeclareError("INITIATOR_IS_ALREADY_EMPLOYEE")

var InitiatorEmployeeHaveNotEnoughRights = ape.DeclareError("INITIATOR_HAVE_NOT_ENOUGH_PERMISSIONS")
