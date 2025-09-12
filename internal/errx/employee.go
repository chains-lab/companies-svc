package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorEmployeeNotFound = ape.DeclareError("EMPLOYEE_NOT_FOUND")

var ErrorEmployeeAlreadyExists = ape.DeclareError("EMPLOYEE_ALREADY_EXISTS")

var EmployeeRoleNotSupported = ape.DeclareError("EMPLOYEE_ROLE_NOT_SUPPORTED")

var ErrorInitiatorNotEmployee = ape.DeclareError("INITIATOR_NOT_EMPLOYEE")

var ErrorInitiatorIsAlreadyEmployee = ape.DeclareError("INITIATOR_IS_ALREADY_EMPLOYEE")

var ErrorInitiatorEmployeeHaveNotEnoughRights = ape.DeclareError("INITIATOR_HAVE_NOT_ENOUGH_PERMISSIONS")

var ErrorCurrentEmployeeCanCreateDistributor = ape.DeclareError("CURRENT_EMPLOYEE_CAN_CREATE_DISTRIBUTOR")

var InitiatorAndUserHaveDifferentDistributors = ape.DeclareError("INITIATOR_AND_USER_HAVE_DIFFERENT_DISTRIBUTORS")
