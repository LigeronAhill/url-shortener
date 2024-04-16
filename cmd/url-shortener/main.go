package main

import (
	"log/slog"
	"os"

	"github.com/LigeronAhill/url-shortener/internal/config"
	"github.com/LigeronAhill/url-shortener/internal/lib/logger/sl"
	"github.com/LigeronAhill/url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config
	cfg := config.MustLoad("config/local.yaml")

	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	g, err := storage.GetURL("googled")
	if err != nil {
		log.Error("failed to get url", sl.Err(err))
		os.Exit(1)
	}
	log.Info("Got url for 'google' alias:", slog.String("url", g))

	// TODO: init router: chi, chi render

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
