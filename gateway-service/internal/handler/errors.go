package handler

import (
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) logError(r *http.Request, err error) {
	log.Print(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := h.encodeJSON(w, status, env, nil)
	if err != nil {
		h.logError(r, err)
		w.WriteHeader(500)
	}
}

func (h *Handler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	h.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (h *Handler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	h.errorResponse(w, r, http.StatusNotFound, message)
}

func (h *Handler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	h.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (h *Handler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *Handler) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	h.errorResponse(w, r, http.StatusForbidden, message)
}

func (h *Handler) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *Handler) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *Handler) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}
