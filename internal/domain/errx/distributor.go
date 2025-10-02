package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorDistributorNotFound = ape.DeclareError("DISTRIBUTOR_NOT_FOUND")

var ErrorDistributorIsBlocked = ape.DeclareError("DISTRIBUTOR_STATUS_BLOCKED")

var ErrorInvalidDistributorStatus = ape.DeclareError("INVALID_DISTRIBUTOR_STATUS")

var ErrorCannotSetDistributorStatusBlocked = ape.DeclareError("CANNOT_SET_DISTRIBUTOR_STATUS_BLOCKED")

var ErrorAnswerToInviteForNotActiveDistributor = ape.DeclareError("ANSWER_TO_INVITE_FOR_NOT_ACTIVE_DISTRIBUTOR")
