package products

import (
	"log"
	"net/http"

	"github.com/ThiwankaS/ecom-go-api-project/internal/json"
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

	err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := []string{"Hello", "World"}

	json.Write(w, http.StatusOK, products)
}
