package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorDistributorBlockNotFound = ape.DeclareError("BLOCK_NOT_FOUND")

var ErrorDistributorHaveAlreadyActiveBlock = ape.DeclareError("DISTRIBUTOR_HAVE_ALREADY_ACTIVE_BLOCK")

var ErrorInvalidDistributorBlockStatus = ape.DeclareError("INVALID_BLOCK_STATUS")
