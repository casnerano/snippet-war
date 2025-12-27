package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/casnerano/snippet-war/internal/service/quiz"
	desc "github.com/casnerano/snippet-war/pkg/api/v1/quiz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	listener net.Listener
	grpc     *grpc.Server
}

func New(addr string) (*Server, error) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := Server{
		listener: listener,
	}

	server.grpc = grpc.NewServer()

	reflection.Register(server.grpc)
	desc.RegisterQuizServer(server.grpc, &quiz.Quiz{})

	return &server, nil
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		slog.Info("Shutting down server...")
		s.grpc.GracefulStop()
	}()

	return s.grpc.Serve(s.listener)
}
