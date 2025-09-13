package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInviteNotFound = ape.DeclareError("INVITE_NOT_FOUND")

var ErrorInvalidInviteToken = ape.DeclareError("INVALID_INVITE_TOKEN")

var ErrorInviteAlreadyAnswered = ape.DeclareError("INVITE_ALREADY_ANSWERED")

var ErrorInviteExpired = ape.DeclareError("INVITE_EXPIRED")

var ErrorInvalidEmployeeRole = ape.DeclareError("INVALID_EMPLOYEE_ROLE")

var ErrorInvalidInviteStatus = ape.DeclareError("INVALID_INVITE_STATUS")

var ErrorInvalidDistributorBlockStatus = ape.DeclareError("INVALID_BLOCK_STATUS")

var ErrorUnexpectedInviteStatus = ape.DeclareError("UNEXPECTED_INVITE_STATUS")
