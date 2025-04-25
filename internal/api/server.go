package api

import (
	"app/internal/config"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      newRoutes(cfg),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		slog.Info("signal caught", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdown <- server.Shutdown(ctx)
	}()

	slog.Info("server started running", "port", cfg.Port)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("error on starting server", "err", err)
	}

	err = <-shutdown
	if err != nil {
		slog.Error("error during server shutdown", "err", err)
	} else {
		slog.Info("server shutdown completed gracefully")
	}
}
