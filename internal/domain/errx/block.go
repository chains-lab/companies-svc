package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorCompanyBlockNotFound = ape.DeclareError("BLOCK_NOT_FOUND")

var ErrorCompanyHaveAlreadyActiveBlock = ape.DeclareError("COMPANY_HAVE_ALREADY_ACTIVE_BLOCK")
