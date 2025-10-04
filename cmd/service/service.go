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

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(
		gracefulStop,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	go func() {
		select {
		//nolint:errcheck
		case <-gracefulStop:
			cancel()
			httpSvc.Stop()
			os.Exit(0)
		case <-ctx.Done():
			return
		}
	}()

	if err := httpSvc.Start(ctx); err != nil {
		panic(err)
	}
}
