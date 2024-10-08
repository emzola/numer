package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	invoicepb "github.com/emzola/numer/invoice-service/proto"
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

// encodeJSON serializes data to JSON and writes the appropriate HTTP status code and headers if necessary.
func (h *Handler) encodeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// decodeJSON de-serializes JSON data into Go types.
func (h *Handler) decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

// Convert gRPC Invoices to HTTP Invoices
func convertInvoices(invoices []*invoicepb.Invoice) []InvoiceHTTP {
	httpInvoices := make([]InvoiceHTTP, len(invoices))
	for i, inv := range invoices {
		httpInvoices[i] = InvoiceHTTP{
			InvoiceID:          inv.Id,
			UserID:             inv.UserId,
			CustomerID:         inv.CustomerId,
			IssueDate:          inv.IssueDate.AsTime(),
			DueDate:            inv.DueDate.AsTime(),
			Currency:           inv.Currency,
			Items:              convertInvoiceItems(inv.Items),
			DiscountPercentage: inv.DiscountPercentage,
			AccountName:        inv.AccountName,
			AccountNumber:      inv.AccountNumber,
			BankName:           inv.BankName,
			RoutingNumber:      inv.RoutingNumber,
			Note:               inv.Note,
		}
	}
	return httpInvoices
}

// Convert gRPC InvoiceItems to HTTP InvoiceItems
func convertInvoiceItems(items []*invoicepb.InvoiceItem) []InvoiceItem {
	httpItems := make([]InvoiceItem, len(items))
	for i, item := range items {
		httpItems[i] = InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}
	return httpItems
}

// ReadString reads a url query param and returns a string
func (h *Handler) ReadString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// ReadString reads a url query param and returns an int
func (h *Handler) ReadInt(qs url.Values, key string, defaultValue int) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return i
}
