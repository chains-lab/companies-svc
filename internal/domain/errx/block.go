package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorDistributorBlockNotFound = ape.DeclareError("BLOCK_NOT_FOUND")

//
//var ErrorNoActiveBlockForDistributor = ape.DeclareError("NO_ACTIVE_BLOCK_FOR_DISTRIBUTOR")

var DistributorHaveAlreadyActiveBlock = ape.DeclareError("DISTRIBUTOR_HAVE_ALREADY_ACTIVE_BLOCK")
