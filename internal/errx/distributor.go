package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorDistributorNotFound = ape.DeclareError("DISTRIBUTOR_NOT_FOUND")

var DistributorStatusBlocked = ape.DeclareError("DISTRIBUTOR_STATUS_BLOCKED")

var UnexpectedDistributorSetStatus = ape.DeclareError("UNEXPECTED_DISTRIBUTOR_SET_STATUS")

var InvalidDistributorStatus = ape.DeclareError("INVALID_DISTRIBUTOR_STATUS")
