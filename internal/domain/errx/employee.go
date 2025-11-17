package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorEmployeeNotFound = ape.DeclareError("EMPLOYEE_NOT_FOUND")

var ErrorUserAlreadyInThisCompany = ape.DeclareError("USER_ALREADY_IN_THIS_COMPANY")

var ErrorCannotDeleteYourself = ape.DeclareError("CANNOT_DELETE_YOURSELF")

var ErrorOwnerCannotRefuseSelf = ape.DeclareError("OWNER_CANNOT_REFUSE_SELF")
