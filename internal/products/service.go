package products

import (
	"context"

	repository "github.com/ThiwankaS/ecom-go-api-project/internal/adapters/postgresql/sqlc"
)

// Service defines the business logic for managing products.
// By using an interface, we can easily create mocks for unit testing
// the handler layer without needing a real database.
type Service interface {
	ListProducts(ctx context.Context) ([]repository.Product, error)
}

// svc is the concrete implementation of the Service interface.
// It is unexported (lowercase) to force use of the NewService constructor.
type svc struct {
	repo repository.Querier
}

// NewService creates a new product service instance.
// It accepts a repository.Querier interface, allowing it to work with
// any store that implements the generated sqlc methods.
func NewService(repo repository.Querier) Service {
	return &svc{
		repo: repo,
	}
}

// ListProducts retrieves all products from the data store.
// It delegates the database operation to the repository layer.
func (s *svc) ListProducts(ctx context.Context) ([]repository.Product, error) {
	// Any business logic (like filtering, caching, or data transformation)
	// should happen here before or after calling the repository.
	return s.repo.ListProducts(ctx)
}
