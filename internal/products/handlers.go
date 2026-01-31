package products

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	service Service
}

// creating a new handler
func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// call the ListProducts service
	// return JSON in http response

	products := []string{"Hello", "World"}

	json.NewEncoder(w).Encode(products)
}
