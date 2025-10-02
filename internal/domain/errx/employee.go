package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorEmployeeNotFound = ape.DeclareError("EMPLOYEE_NOT_FOUND")

var ErrorEmployeeAlreadyExists = ape.DeclareError("EMPLOYEE_ALREADY_EXISTS")

var EmployeeInvalidRole = ape.DeclareError("INVALID_EMPLOYEE_ROLE")

var ErrorInitiatorNotEmployee = ape.DeclareError("INITIATOR_NOT_EMPLOYEE")

var ErrorInitiatorIsAlreadyEmployee = ape.DeclareError("INITIATOR_IS_ALREADY_EMPLOYEE")

var ErrorCannotDeleteYourself = ape.DeclareError("CANNOT_DELETE_YOURSELF")

var ErrorInitiatorEmployeeHaveNotEnoughRights = ape.DeclareError("INITIATOR_HAVE_NOT_ENOUGH_PERMISSIONS")

var ErrorInitiatorAndUserHaveDifferentDistributors = ape.DeclareError("INITIATOR_AND_USER_HAVE_DIFFERENT_DISTRIBUTORS")

var ErrorInitiatorIsNotThisDistributorEmployee = ape.DeclareError("INITIATOR_IS_NOT_THIS_DISTRIBUTOR_EMPLOYEE")

var ErrorCurrentEmployeeCannotCreateDistributor = ape.DeclareError("CURRENT_EMPLOYEE_CANNOT_CREATE_DISTRIBUTOR")
