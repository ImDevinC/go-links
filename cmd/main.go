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

	ctx := context.Background()
	cfg, err := config.FromEnv(ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	if err := server.Start(ctx, &cfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
