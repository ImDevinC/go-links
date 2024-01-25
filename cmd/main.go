package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/imdevinc/go-links/internal/app"
	"github.com/imdevinc/go-links/internal/config"
	"github.com/imdevinc/go-links/internal/store"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := app.App{
		Store:  store.NewMemoryStore(),
		Logger: logger,
	}

	server.Store.CreateLink(context.Background(), store.Link{
		Name:        "site",
		URL:         "https://imdevinc.com",
		Description: "site",
		Views:       2,
	})
	server.Store.CreateLink(context.Background(), store.Link{
		Name:        "imdevinc",
		URL:         "https://github.com/imdevinc",
		Description: "GitHub",
		Views:       10,
	})
	server.Store.CreateLink(context.Background(), store.Link{
		Name:        "google",
		URL:         "https://google.com",
		Description: "google",
		Views:       1,
	})

	cfg := config.FromEnv()
	err := server.Start(context.Background(), &cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
