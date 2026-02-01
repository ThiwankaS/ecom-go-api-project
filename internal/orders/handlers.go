package orders

import (
	"log"
	"net/http"

	"github.com/ThiwankaS/ecom-go-api-project/internal/json"
)

// handler represents the HTTP dependencies for product-related endpoints.
// Using an unexported struct helps encapsulate the routing logic.
type handler struct {
	service Service
}

// NewHandler initializes a new products handler with the provided service.
// This follows the dependency injection pattern, making it easier to test.
func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// PlaceOrder handles GET requests for the product catalog.
// It retrieves products from the service and returns them as a JSON array.
func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder CreateOrderItemParams
	if err := json.Read(r, &tempOrder); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// internal/json helper ensures consistent response formatting
	json.Write(w, http.StatusCreated, createdOrder)
}
