package main

import (
	"GoDojoGo/deff"
	globals "GoDojoGo/deff"
	"GoDojoGo/libbrevo"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cron "GoDojoGo/worker"

	"github.com/labstack/echo/v5"
)

type Template struct {
	t *template.Template
}

// Echo v5 Renderer interface
func (r *Template) Render(c *echo.Context, w io.Writer, name string, data any) error {
	return r.t.ExecuteTemplate(w, name, data)
}

func RunWithWorker() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start worker
	go cron.Run(ctx, deff.Settings.POOL_CHECK)

	e := InitService()

	// Create HTTP server manually
	server := &http.Server{
		Addr:    ":" + fmt.Sprint(globals.Settings.PORT),
		Handler: e,
	}

	// Start HTTP server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.Logger.Error("failed to start server", "error", err)
		}
	}()

	// Wait for signal
	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		e.Logger.Error("server shutdown failed", "error", err)
	}
}

func RunWithoutWorker() {

	e := InitService()

	// Create HTTP server manually
	server := &http.Server{
		Addr:    ":" + fmt.Sprint(globals.Settings.PORT),
		Handler: e,
	}

	// Start HTTP server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.Logger.Error("failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		e.Logger.Error("server shutdown failed", "error", err)
	}
}

func main() {

	runWorker := true

	if err := globals.LoadSettings(".env"); err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "version":
			fmt.Println(deff.VERSION)
			return
		case "noworker":
			runWorker = false
		case "brevocheck":
			rep, err := libbrevo.SendinblueCheck(deff.Settings.BREVO_API_KEY)
			if err == nil {
				for _, r := range rep {
					fmt.Println(r)
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
	}

	if runWorker {
		fmt.Println("Running with worker...")
		RunWithWorker()
	} else {
		fmt.Println("Running without worker...")
		RunWithoutWorker()
	}

}
