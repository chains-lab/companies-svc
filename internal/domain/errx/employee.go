package errx

import (
	"github.com/chains-lab/ape"
)

var EmployeeInvalidRole = ape.DeclareError("INVALID_EMPLOYEE_ROLE")

var ErrorEmployeeNotFound = ape.DeclareError("EMPLOYEE_NOT_FOUND")

var ErrorEmployeeAlreadyExists = ape.DeclareError("EMPLOYEE_ALREADY_EXISTS")

var ErrorCannotDeleteYourself = ape.DeclareError("CANNOT_DELETE_YOURSELF")

var ErrorOwnerCannotRefuseSelf = ape.DeclareError("OWNER_CANNOT_REFUSE_SELF")

var ErrorInitiatorHaveNotEnoughRights = ape.DeclareError("INITIATOR_HAVE_NOT_ENOUGH_PERMISSIONS")

var ErrorInitiatorIsNotEmployee = ape.DeclareError("INITIATOR_NOT_EMPLOYEE")

var ErrorInitiatorIsNotEmployeeOfThisCompany = ape.DeclareError("INITIATOR_NOT_EMPLOYEE_OF_THIS_company")
