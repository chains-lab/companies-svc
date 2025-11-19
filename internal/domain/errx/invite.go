package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInviteNotFound = ape.DeclareError("INVITE_NOT_FOUND")

var ErrorInviteAlreadyReplyed = ape.DeclareError("INVITE_ALREADY_REPLYED")

var ErrorInviteExpired = ape.DeclareError("INVITE_EXPIRED")

var ErrorInvalidEmployeeRole = ape.DeclareError("INVALID_EMPLOYEE_ROLE")

var ErrorInvalidInviteStatus = ape.DeclareError("INVALID_INVITE_STATUS")

var ErrorInviteNotForThisUser = ape.DeclareError("INVITE_NOT_FOR_THIS_USER")
