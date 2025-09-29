package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tsetsik/vehicle-search/internal/http"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	httpSvc, err := http.NewService()
	if err != nil {
		panic(err)
	}

	if err := httpSvc.Start(ctx); err != nil {
		panic(err)
	}

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(
		gracefulStop,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	go func() {
		select {
		case <-gracefulStop:
		case <-ctx.Done():
			cancel()
			//nolint:errcheck
			httpSvc.Stop()
			// handle it
			os.Exit(0)
		}

	}()
}
