package orders

import (
	"log"
	"net/http"

	"github.com/ThiwankaS/ecom-go-api-project/internal/json"
)

// handler coordinates the HTTP layer for order processing.
// It acts as the primary entry point for customers to submit new purchases.
type handler struct {
	service Service
}

// NewHandler creates a new instance of the order handler.
// Dependency injection of the Service allows for cleaner testing of business rules.
func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// PlaceOrder processes an incoming purchase request.
// It decodes the request body, validates the items, and triggers the order workflow.
func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder CreateOrderItemParams
	// Parse the JSON request body into our parameter struct.
	// If decoding fails, we return a 400 Bad Request to indicate a client error.
	if err := json.Read(r, &tempOrder); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Pass the request context and parameters to the service layer.
	// The service handles complex logic like stock checking and payment.
	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Println(err)
		// Generic 500 error prevents leaking internal DB details to the client.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created order with a 201 Created status.
	// The response body typically includes the generated Order ID and total price.
	json.Write(w, http.StatusCreated, createdOrder)
}
