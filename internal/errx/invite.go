package errx

import (
	"github.com/chains-lab/ape"
)

var InviteNotFound = ape.DeclareError("INVITE_NOT_FOUND")

var InviteIsNotActive = ape.DeclareError("INVITE_IS_NOT_ACTIVE")

var InviteIsNotForInitiator = ape.DeclareError("INVITE_IS_NOT_FOR_INITIATOR")

var CantSendInviteForCurrentEmployee = ape.DeclareError("CANT_SEND_INVITE_FOR_CURRENT_EMPLOYEE")

var UserHaveAlreadyInviteForInitiatorDistributor = ape.DeclareError("USER_HAVE_ALREADY_INVITE_FOR_INITIATOR_DISTRIBUTOR")
