package orders

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/ThiwankaS/ecom-go-api-project/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

// Define custom domain errors to separate database technicalities
// from clear business logic failures.
var (
	ErrProductNotFound    = errors.New("Product Not Found")
	ErrProductOutOfStock  = errors.New("Product Out Of Stock")
	ErrProductStockUpdate = errors.New("Product Stock Update failed")
)

// svc manages the lifecycle of an order.
// It requires both the generated Queries and a raw pgx.Conn to handle transactions.
type svc struct {
	repo *repository.Queries
	db   *pgx.Conn
}

// NewService initializes the order business logic layer.
func NewService(repo *repository.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

// PlaceOrder executes the full checkout workflow.
// This function uses a Database Transaction to ensure that if one step fails
// (like a product being out of stock), no partial order data is saved.
func (s *svc) PlaceOrder(ctx context.Context, tempOrder CreateOrderItemParams) (repository.Order, error) {
	// 1. Basic Request Validation
	// We check these before opening a DB connection to save resources.
	if tempOrder.CustomerID == 0 {
		return repository.Order{}, fmt.Errorf("Customer ID required")
	}

	if len(tempOrder.Items) == 0 {
		return repository.Order{}, fmt.Errorf("At least one item required")
	}

	// 2. Start Transaction
	// All operations from here on are "All or Nothing."
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repository.Order{}, err
	}
	// Defer a Rollback. If the function returns early due to an error,
	// tx.Rollback ensures no garbage data remains in the DB.
	defer tx.Rollback(ctx)
	// Attach the transaction to our sqlc repository.
	qtx := s.repo.WithTx(tx)

	// 3. Create the parent Order record
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repository.Order{}, err
	}

	// 4. Process each Line Item
	for _, item := range tempOrder.Items {
		// Verify product existence
		product, err := qtx.ListProductById(ctx, item.ProductID)
		if err != nil {
			return repository.Order{}, ErrProductNotFound
		}
		// Inventory Check
		if product.Quantity < item.Quantity {
			return repository.Order{}, ErrProductOutOfStock
		}
		// Save the order item
		// We capture the price at the moment of purchase to protect against future price changes.
		_, err = qtx.CreateOrderItem(ctx, repository.CreateOrderItemParams{
			OrderID:      order.ID,
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PriceInCents: product.PriceInCents,
		})

		if err != nil {
			return repository.Order{}, err
		}

		// Deduct Stock
		// We do this AFTER creating the order item but BEFORE the commit
		err = qtx.UpdateProductStock(ctx, repository.UpdateProductStockParams{
			Quantity: item.Quantity,
			ID:       item.ProductID,
		})

		if err != nil {
			return repository.Order{}, ErrProductStockUpdate
		}
	}
	// 5. Commit Transaction
	// If we reached here, every item was valid and in stock.
	if err := tx.Commit(ctx); err != nil {
		return repository.Order{}, err
	}

	return order, nil
}
