package app

import (
	"GoLoad/internal/handler/consumers"
	"GoLoad/internal/handler/grpc"
	"GoLoad/internal/handler/http"
	"GoLoad/internal/utils"
	"context"
	"go.uber.org/zap"
	"syscall"
)

type Server struct {
	grpcServer   grpc.Server
	httpServer   http.Server
	rootConsumer consumers.Root
	logger       *zap.Logger
}

func NewServer(
	grpcServer grpc.Server,
	httpServer http.Server,
	rootConsumer consumers.Root,
	logger *zap.Logger,
) *Server {
	return &Server{
		grpcServer:   grpcServer,
		httpServer:   httpServer,
		rootConsumer: rootConsumer,
		logger:       logger,
	}
}

func (s Server) Start() error {
	go func() {
		err := s.grpcServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("gRPC server stopped")
	}()

	go func() {
		err := s.httpServer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("HTTP server stopped")
	}()

	go func() {
		err := s.rootConsumer.Start(context.Background())
		s.logger.With(zap.Error(err)).Info("Message queue consumer stopped")
	}()

	utils.BlockUntilSignal(syscall.SIGINT, syscall.SIGINT)
	return nil
}
