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
	ctx := context.Background()

	cfg, err := config.FromEnv(ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var s store.Store
	switch cfg.StoreType {
	case config.StoreTypeFile:
		s, err = store.NewFileStore("links.json", true)
	case config.StoreTypeMemory:
		s = store.NewMemoryStore()
	case config.StoreTypeMongo:
		s, err = store.NewMongoDBStore(ctx, cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.DatabaseName)
	}

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer s.Close(ctx)

	server := app.App{
		Store:  s,
		Logger: logger,
	}

	if err := server.Start(ctx, &cfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
