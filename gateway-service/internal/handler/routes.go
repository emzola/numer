package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
)

func (h *Handler) Routes(userServiceConn *grpc.ClientConn) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(h.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(h.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/invoices", h.authMiddleware(h.GetInvoicesHandler, userServiceConn))
	router.HandlerFunc(http.MethodPost, "/invoices", h.authMiddleware(h.CreateInvoiceHandler, userServiceConn))
	router.HandlerFunc(http.MethodGet, "/invoices/:id", h.authMiddleware(h.GetInvoiceHandler, userServiceConn))
	router.HandlerFunc(http.MethodPatch, "/invoices/:id", h.authMiddleware(h.UpdateInvoiceHandler, userServiceConn))
	router.HandlerFunc(http.MethodPost, "/invoices/:id/send", h.authMiddleware(h.SendInvoiceHandler, userServiceConn))

	router.HandlerFunc(http.MethodGet, "/stats", h.authMiddleware(h.GetStatsHandler, userServiceConn))

	router.HandlerFunc(http.MethodGet, "/activities/user", h.authMiddleware(h.GetUserActivitiesHandler, userServiceConn))
	router.HandlerFunc(http.MethodGet, "/activities/invoice/:id", h.authMiddleware(h.GetInvoiceActivitiesHandler, userServiceConn))

	router.HandlerFunc(http.MethodPost, "/reminders", h.authMiddleware(h.GetInvoiceActivitiesHandler, userServiceConn))

	router.HandlerFunc(http.MethodPost, "/users", h.CreateUserHandler)
	router.HandlerFunc(http.MethodGet, "/users/:id", h.authMiddleware(h.GetUserHandler, userServiceConn))
	router.HandlerFunc(http.MethodPatch, "/users/:id", h.authMiddleware(h.UpdateUserHandler, userServiceConn))
	router.HandlerFunc(http.MethodDelete, "/users/:id", h.authMiddleware(h.DeleteUserHandler, userServiceConn))

	router.HandlerFunc(http.MethodPost, "/customers", h.authMiddleware(h.CreateCustomerHandler, userServiceConn))
	router.HandlerFunc(http.MethodGet, "/customers/:id", h.authMiddleware(h.GetCustomerHandler, userServiceConn))
	router.HandlerFunc(http.MethodPatch, "/customers/:id", h.authMiddleware(h.UpdateCustomerHandler, userServiceConn))
	router.HandlerFunc(http.MethodDelete, "/customers/:id", h.authMiddleware(h.DeleteCustomerHandler, userServiceConn))

	return router
}
