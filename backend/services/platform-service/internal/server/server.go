package server

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serviceRegistrar interface {
	Register(server *grpc.Server)
}

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
		grpc:     grpc.NewServer(),
	}

	return &server, nil
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		slog.Info("Shutting down server...")

		s.grpc.GracefulStop()
		_ = s.listener.Close()
	}()

	return s.grpc.Serve(s.listener)
}

func (s *Server) RegisterServices(services ...serviceRegistrar) {
	for _, service := range services {
		service.Register(s.grpc)
	}
	reflection.Register(s.grpc)
}
