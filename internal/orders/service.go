package orders

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/ThiwankaS/ecom-go-api-project/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound   = errors.New("Product Not Found")
	ErrProductOutOfStock = errors.New("Product Out Of Stock")
)

// svc is the concrete implementation of the Service interface.
// It is unexported (lowercase) to force use of the NewService constructor.
type svc struct {
	repo *repository.Queries
	db   *pgx.Conn
}

// NewService creates a new product service instance.
// It accepts a repository.Querier interface, allowing it to work with
// any store that implements the generated sqlc methods.
func NewService(repo *repository.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

// PlaceOrder retrieves all products from the data store.
// It delegates the database operation to the repository layer.
func (s *svc) PlaceOrder(ctx context.Context, tempOrder CreateOrderItemParams) (repository.Order, error) {
	// Any business logic (like filtering, caching, or data transformation)
	// should happen here before or after calling the repository.

	// Incoming data validation
	if tempOrder.CustomerID == 0 {
		return repository.Order{}, fmt.Errorf("Customer ID required")
	}

	if len(tempOrder.Items) == 0 {
		return repository.Order{}, fmt.Errorf("At least one item required")
	}

	// If item is not existing on the DB let's create on
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repository.Order{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// create an order
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repository.Order{}, err
	}

	// check all the productID one by for existance
	for _, item := range tempOrder.Items {

		product, err := qtx.ListProductById(ctx, item.ProductID)
		if err != nil {
			return repository.Order{}, ErrProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repository.Order{}, ErrProductOutOfStock
		}

		_, err = qtx.CreateOrderItem(ctx, repository.CreateOrderItemParams{
			OrderID:      order.ID,
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PriceInCents: product.PriceInCenters,
		})

		if err != nil {
			return repository.Order{}, err
		}
	}

	tx.Commit(ctx)

	return order, nil
}
