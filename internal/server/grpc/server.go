package grpc

import (
	"context"
	"fmt"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server/grpc/interceptors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	s      *grpc.Server
	port   string
	logger *zap.Logger
}

func NewServer(
	lc fx.Lifecycle,
	loginServiceServer *LoginServiceServer,
	noteServiceServer *NoteServiceServer,
	binaryServiceServer *BinaryServiceServer,
	cardServiceServer *CardServiceServer,
	cfg *Config,
	logger *zap.Logger,
) *Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.UnaryLogger(logger),
			interceptors.UnaryAuth(),
		),
		grpc.ChainStreamInterceptor(
			interceptors.StreamLogger(logger),
			interceptors.StreamAuth(),
		),
	)
	pb.RegisterLoginServiceServer(s, loginServiceServer)
	pb.RegisterNoteServiceServer(s, noteServiceServer)
	pb.RegisterBinaryServiceServer(s, binaryServiceServer)
	pb.RegisterCardServiceServer(s, cardServiceServer)
	reflection.Register(s)

	srv := &Server{
		s:      s,
		port:   cfg.Port,
		logger: logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go srv.start()
			return nil
		},
		OnStop: func(context.Context) error {
			srv.s.GracefulStop()
			return nil
		},
	})

	return srv
}

func (s *Server) start() error {
	s.logger.Info("running grpc server", zap.String("port", s.port))

	listen, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := s.s.Serve(listen); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	s.logger.Info("server stopped")

	return nil
}
