package main

import (
	"github.com/Slava02/URL_shortner/internal/config"
	"github.com/Slava02/URL_shortner/internal/lib/logger/sl"
	"github.com/Slava02/URL_shortner/internal/storage/sqlite"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//  TODO: init config: cleanenv
	cfg := config.MustLoad()

	//  TODO: init logger: sl
	log := setUpLogger(cfg.Env)

	log.Info("starting url-shortner", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	//  TODO: init db: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
	//  TODO: init router: chi, chi render

	//  TODO: run server
}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
