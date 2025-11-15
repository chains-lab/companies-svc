package requests

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/companies-svc/resources"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UpdateEmployee(r *http.Request) (req resources.UpdateEmployee, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.EmployeeType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "user_id") != req.Data.Id.String() {
		errs["data/id"] = validation.NewError("mismatch", "id in path and body must match")
	}

	return req, errs.Filter()
}
