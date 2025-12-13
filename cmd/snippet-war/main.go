package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/casnerano/snippet-war/internal/server"
)

const defaultAddr = ":8086"

var addr = flag.String("addr", defaultAddr, "Server address")

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	flag.Parse()

	srv := server.New(*addr)

	if err := srv.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server failed: %v", err)
	}
}
