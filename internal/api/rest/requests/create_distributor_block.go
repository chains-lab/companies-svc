package requests

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/distributors-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateDistributorBlock(r *http.Request) (req resources.CreateDistributorBlock, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.CreateDistributorBlockType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
	return req, errs.Filter()
}
