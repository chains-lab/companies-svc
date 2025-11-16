package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorEmployeeNotFound = ape.DeclareError("EMPLOYEE_NOT_FOUND")

var ErrorUserAlreadyEmployee = ape.DeclareError("USER_ALREADY_EMPLOYEE")

var ErrorCannotDeleteYourself = ape.DeclareError("CANNOT_DELETE_YOURSELF")

var ErrorOwnerCannotRefuseSelf = ape.DeclareError("OWNER_CANNOT_REFUSE_SELF")

var ErrorInitiatorIsNotEmployee = ape.DeclareError("INITIATOR_IS_NOT_EMPLOYEE")

var ErrorInitiatorIsNotEmployeeInThisCompany = ape.DeclareError("INITIATOR_IS_IN_NOT_EMPLOYEE_THIS_COMPANY")
