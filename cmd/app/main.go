package main

import (
	"app/internal/api"
	"app/internal/config"
	"log/slog"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	slog.Info("env parsed successfully", "environment", cfg.Env)

	api.Run(cfg)
}
