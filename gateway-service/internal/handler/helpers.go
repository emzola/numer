package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// envelop is a wrapper around JSON responses.
type envelope map[string]interface{}

// readIDParam pulls the url id parameter from the request and returns it or an error if any
func (h *Handler) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
