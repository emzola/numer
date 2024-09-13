package handler

import (
	"context"
	"net/http"
	"time"

	activitypb "github.com/emzola/numer/activity-service/proto"
	"github.com/emzola/numer/gateway-service/internal/grpcutil"
)

func (h *Handler) GetUserActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	user := h.contextGetUser(r)

	// Extract ID param
	userId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
	}

	// Basic validation to ensure user fetches own activities
	if user.Id != userId {
		h.notPermittedResponse(w, r)
		return
	}

	// Convert the HTTP request into the gRPC ListInvoicesRequest
	grpcReq := &activitypb.GetUserActivitiesRequest{
		UserId: user.Id,
		Limit:  10,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to activity service
	conn, err := grpcutil.ServiceConnection(ctx, "activity-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := activitypb.NewActivityServiceClient(conn)
	grpcRes, err := client.GetUserActivities(context.Background(), grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC GetUserActivitiesResponse to the HTTP response
	var activities []ActivityHTTPResp
	for _, activity := range grpcRes.Activities {
		activities = append(activities, ActivityHTTPResp{
			InvoiceID:   activity.InvoiceId,
			UserID:      activity.UserId,
			Action:      activity.Action,
			Description: activity.Description,
			Timestamp:   activity.Timestamp,
		})
	}

	userActResp := GetUserActivitiesHTTPResp{
		Activities: activities,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"user activities": userActResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) GetInvoiceActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID param
	invoiceId, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	// Convert the HTTP request into the gRPC ListInvoicesRequest
	grpcReq := &activitypb.GetInvoiceActivitiesRequest{
		InvoiceId: invoiceId,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Create gRPC connection to activity service
	conn, err := grpcutil.ServiceConnection(ctx, "activity-service", h.registry)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
	defer conn.Close()

	client := activitypb.NewActivityServiceClient(conn)
	grpcRes, err := client.GetInvoiceActivities(context.Background(), grpcReq)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	// Map the gRPC GetInvoiceActivitiesResponse to the HTTP response
	var activities []ActivityHTTPResp
	for _, activity := range grpcRes.Activities {
		activities = append(activities, ActivityHTTPResp{
			InvoiceID:   activity.InvoiceId,
			UserID:      activity.UserId,
			Action:      activity.Action,
			Description: activity.Description,
			Timestamp:   activity.Timestamp,
		})
	}

	invActResp := GetInvoiceActivitiesHTTPResp{
		Activities: activities,
	}

	err = h.encodeJSON(w, http.StatusOK, envelope{"invoice activities": invActResp}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// Struct to capture the HTTP response
type GetUserActivitiesHTTPResp struct {
	Activities []ActivityHTTPResp `json:"activities"`
}

// Struct to capture each Activity in the HTTP response
type ActivityHTTPResp struct {
	InvoiceID   int64  `json:"invoice_id"`
	UserID      int64  `json:"user_id"`
	Action      string `json:"action"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// Struct to capture the HTTP response (matching GetInvoiceActivitiesResponse)
type GetInvoiceActivitiesHTTPResp struct {
	Activities []ActivityHTTPResp `json:"activities"`
}
