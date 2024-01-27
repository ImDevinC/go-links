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

	cfg := config.FromEnv()
	err := server.Start(context.Background(), &cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
