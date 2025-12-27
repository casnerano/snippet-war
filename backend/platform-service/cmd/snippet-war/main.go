package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/casnerano/snippet-war/internal/server"
)

const defaultAddr = ":8081"

var addr string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	flag.StringVar(&addr, "addr", defaultAddr, "Server address")
	flag.Parse()

	// Создание сервера с сервисом
	srv := server.New(addr)

	if err := srv.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		os.Exit(1)
	}
}
