package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/imdevinc/go-links/internal/app"
	"github.com/imdevinc/go-links/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := app.App{
		Store:  store.NewMemoryStore(),
		Logger: logger,
	}

	err := server.Store.CreateLink(context.Background(), store.Link{
		Name:        "link",
		URL:         "https://google.com",
		Description: "google",
	})

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = server.Start(context.Background())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
