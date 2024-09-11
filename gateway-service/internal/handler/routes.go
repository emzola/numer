package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) Routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/invoices", h.GetInvoices)
	router.HandlerFunc(http.MethodPost, "/invoices", h.CreateInvoice)
	router.HandlerFunc(http.MethodGet, "/invoices/:id", h.GetInvoice)
	router.HandlerFunc(http.MethodPatch, "/invoices/:id", h.UpdateInvoice)

	return router
}
