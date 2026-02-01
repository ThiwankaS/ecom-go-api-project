package orders

import (
	"context"

	repository "github.com/ThiwankaS/ecom-go-api-project/internal/adapters/postgresql/sqlc"
)

type OrderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type CreateOrderItemParams struct {
	CustomerID int64       `json:"customerId"`
	Items      []OrderItem `json:"items"`
}

type Service interface {
	PlaceOrder(ctx context.Context, tempOrder CreateOrderItemParams) (repository.Order, error)
}
