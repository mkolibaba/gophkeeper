package grpc

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"github.com/mkolibaba/gophkeeper/internal/server/grpc/interceptors"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	s      *grpc.Server
	port   string
	logger *log.Logger
}

type ServerParams struct {
	fx.In

	Lifecycle                  fx.Lifecycle
	AuthService                *server.AuthService
	AuthorizationServiceServer *AuthorizationServiceServer
	LoginServiceServer         *LoginServiceServer
	NoteServiceServer          *NoteServiceServer
	BinaryServiceServer        *BinaryServiceServer
	CardServiceServer          *CardServiceServer
	Config                     *Config
	Logger                     *log.Logger
}

func NewServer(p ServerParams) *Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.UnaryLogger(p.Logger),
			interceptors.UnaryAuth(p.AuthService),
		),
		grpc.ChainStreamInterceptor(
			interceptors.StreamLogger(p.Logger),
			interceptors.StreamAuth(p.AuthService),
		),
	)
	pb.RegisterAuthorizationServiceServer(s, p.AuthorizationServiceServer)
	pb.RegisterLoginServiceServer(s, p.LoginServiceServer)
	pb.RegisterNoteServiceServer(s, p.NoteServiceServer)
	pb.RegisterBinaryServiceServer(s, p.BinaryServiceServer)
	pb.RegisterCardServiceServer(s, p.CardServiceServer)
	reflection.Register(s)

	srv := &Server{
		s:      s,
		port:   p.Config.Port,
		logger: p.Logger,
	}

	p.Lifecycle.Append(fx.Hook{
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
	s.logger.Info("running grpc server", "port", s.port)

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
