package main

import (
	"GoDojoGo/deff"
	worker "GoDojoGo/worker"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runHTTPServer(ctx context.Context, e http.Handler, log interface {
	Error(msg string, keysAndValues ...interface{})
}) {
	server := &http.Server{
		Addr:    ":" + fmt.Sprint(deff.Settings.PORT),
		Handler: e,
	}

	// Start HTTP server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", "error", err)
		}
	}()

	// Wait for shutdown signal/context cancellation
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown failed", "error", err)
	}
}

func RunWithWorker() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start worker
	go worker.Run(ctx, deff.Settings.POOL_CHECK)

	e := InitService()

	// echo.Echo has Logger with Error(...)
	runHTTPServer(ctx, e, e.Logger)
}

func RunWithoutWorker() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	e := InitService()
	runHTTPServer(ctx, e, e.Logger)
}
