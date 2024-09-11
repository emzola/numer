package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) Routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/invoices", h.GetInvoices)
	router.HandlerFunc(http.MethodGet, "/invoices{id}", h.GetInvoices)

	return router
}
