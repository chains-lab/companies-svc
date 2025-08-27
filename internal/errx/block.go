package errx

import (
	"github.com/chains-lab/ape"
)

var DistributorBlockNotFound = ape.DeclareError("BLOCK_NOT_FOUND")

var DistributorHaveAlreadyActiveBlock = ape.DeclareError("DISTRIBUTOR_HAVE_ALREADY_ACTIVE_BLOCK")
