package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var ErrorNotEnoughRight = ape.DeclareError("ERROR_NOT_ENOUGH_RIGHT")
