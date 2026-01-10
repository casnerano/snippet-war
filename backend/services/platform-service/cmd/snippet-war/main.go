package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/casnerano/snippet-war/internal/bootstrap"
	"github.com/casnerano/snippet-war/internal/config"
	"github.com/casnerano/snippet-war/internal/server"
)

const defaultAddr = ":8081"

var addr string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := config.Init(ctx); err != nil {
		log.Fatal("failed to initialize config", err)
	}

	flag.StringVar(&addr, "addr", defaultAddr, "Server address")
	flag.Parse()

	grpcServer, err := server.New(addr)
	if err != nil {
		log.Fatal("failed initialization grpc server: ", err.Error())
	}

	servers, err := bootstrap.InitServers(ctx)
	if err != nil {
		log.Fatal("failed initialization servers: ", err.Error())
	}

	grpcServer.RegisterServices(
		servers.Quiz,
	)

	if err = grpcServer.Run(ctx); err != nil {
		log.Fatal("failed run grpc server: ", err.Error())
	}
}
