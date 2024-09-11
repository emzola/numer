package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (h *Handler) Routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/invoices", h.GetInvoicesHandler)
	router.HandlerFunc(http.MethodPost, "/invoices", h.CreateInvoiceHandler)
	router.HandlerFunc(http.MethodGet, "/invoices/:id", h.GetInvoiceHandler)
	router.HandlerFunc(http.MethodPatch, "/invoices/:id", h.UpdateInvoiceHandler)

	router.HandlerFunc(http.MethodGet, "/stats", h.GetStatsHandler)

	router.HandlerFunc(http.MethodPost, "/users", h.CreateUserHandler)
	router.HandlerFunc(http.MethodGet, "/users/:id", h.GetUserHandler)
	router.HandlerFunc(http.MethodPatch, "/users/:id", h.UpdateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/users/:id", h.DeleteUserHandler)

	router.HandlerFunc(http.MethodPost, "/customers", h.CreateCustomerHandler)
	router.HandlerFunc(http.MethodGet, "/customers/:id", h.GetCustomerHandler)
	router.HandlerFunc(http.MethodPatch, "/customers/:id", h.UpdateCustomerHandler)
	router.HandlerFunc(http.MethodDelete, "/customers/:id", h.DeleteCustomerHandler)

	return router
}
