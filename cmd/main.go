package main

import (
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}

	api := application{
		config: cfg,
	}

	// structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server failed to start %s", "error", err)
		os.Exit(1)
	}
}
