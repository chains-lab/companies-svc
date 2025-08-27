package errx

import (
	"github.com/chains-lab/ape"
)

var DistributorNotFound = ape.DeclareError("DISTRIBUTOR_NOT_FOUND")

var DistributorStatusBlocked = ape.DeclareError("DISTRIBUTOR_STATUS_BLOCKED")
