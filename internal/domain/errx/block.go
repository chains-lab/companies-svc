package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorcompanyBlockNotFound = ape.DeclareError("BLOCK_NOT_FOUND")

var ErrorcompanyHaveAlreadyActiveBlock = ape.DeclareError("company_HAVE_ALREADY_ACTIVE_BLOCK")

var ErrorInvalidCompanyBlockStatus = ape.DeclareError("INVALID_BLOCK_STATUS")
