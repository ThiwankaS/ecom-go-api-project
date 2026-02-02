package main

import (
	"log"
	"net/http"
	"time"

	repository "github.com/ThiwankaS/ecom-go-api-project/internal/adapters/postgresql/sqlc"
	"github.com/ThiwankaS/ecom-go-api-project/internal/orders"
	"github.com/ThiwankaS/ecom-go-api-project/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string // address of the server
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good...\n"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good...\n"))
	})

	productService := products.NewService(repository.New(app.db))
	prodcutHandler := products.NewHandler(productService)
	r.Get("/products", prodcutHandler.ListProducts)

	orderService := orders.NewService(repository.New(app.db), app.db)
	orderHandler := orders.NewHandler(orderService);
	r.Post("/orders", orderHandler.PlaceOrder)

	return r
}

// run
func (app *application) run(h http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server started to listen on address %s", app.config.addr)

	return srv.ListenAndServe()
}
