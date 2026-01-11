package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	app_config "github.com/casnerano/snippet-war/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	content_client "github.com/casnerano/snippet-war/internal/client/content_service"
	quiz_handler "github.com/casnerano/snippet-war/internal/handler/quiz"
	quiz_service "github.com/casnerano/snippet-war/internal/service/quiz"
	quiz_desc "github.com/casnerano/snippet-war/pkg/api/v1/quiz"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	config, err := app_config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err)
	}

	listener, err := net.Listen("tcp", config.Server.GRPC.Addr)
	if err != nil {
		log.Fatalf("Failed to listen %s: %s\n", config.Server.GRPC.Addr, err)
	}
	defer func() {
		_ = listener.Close()
	}()

	grpcServer := grpc.NewServer()

	contentServiceClient := getContentServiceClient(ctx, config.ContentService.Addr)
	quizHandler := getQuizHandler(contentServiceClient)

	quiz_desc.RegisterQuizServer(grpcServer, quizHandler)

	reflection.Register(grpcServer)

	gwMux := runtime.NewServeMux()

	mux := http.NewServeMux()
	mux.Handle("/api/", gwMux)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err = quiz_desc.RegisterQuizHandlerFromEndpoint(ctx, gwMux, config.Server.GRPC.Addr, opts)
	if err != nil {
		log.Fatalf("Failed to register quiz handler: %s\n", err)
	}

	httpServer := &http.Server{
		Addr:    config.Server.HTTP.Addr,
		Handler: mux,
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		log.Printf("Starting gRPC server at %s\n", config.Server.GRPC.Addr)

		defer wg.Done()
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC at %s: %s\n", config.Server.GRPC.Addr, err)
		}
	}()

	wg.Add(1)
	go func() {
		log.Printf("Starting HTTP server at %s\n", config.Server.HTTP.Addr)

		defer wg.Done()
		if err = httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to serve HTTP at %s: %s\n", config.Server.HTTP.Addr, err)
		}
	}()

	<-ctx.Done()

	slog.Info("Shutting down server...")

	_ = httpServer.Shutdown(ctx)
	grpcServer.GracefulStop()
	_ = listener.Close()
}

func getQuizHandler(contentClient *content_client.Client) *quiz_handler.Quiz {
	quizService := quiz_service.New(contentClient)
	return quiz_handler.NewQuiz(quizService)
}

func getContentServiceClient(ctx context.Context, addr string) *content_client.Client {
	return content_client.New(ctx, addr)
}
