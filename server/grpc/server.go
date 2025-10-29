package grpc

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/grpc/interceptors"
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
	AuthInterceptor            *interceptors.AuthInterceptor
	LoggerInterceptor          *interceptors.LoggerInterceptor
	AuthorizationServiceServer *AuthorizationServiceServer
	LoginServiceServer         *LoginServiceServer
	NoteServiceServer          *NoteServiceServer
	BinaryServiceServer        *BinaryServiceServer
	CardServiceServer          *CardServiceServer
	Config                     *server.Config
	Logger                     *log.Logger
}

func NewServer(p ServerParams) *Server {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			p.LoggerInterceptor.Unary,
			p.AuthInterceptor.Unary,
		),
		grpc.ChainStreamInterceptor(
			p.LoggerInterceptor.Stream,
			p.AuthInterceptor.Stream,
		),
	)
	gophkeeperv1.RegisterAuthorizationServiceServer(s, p.AuthorizationServiceServer)
	gophkeeperv1.RegisterLoginServiceServer(s, p.LoginServiceServer)
	gophkeeperv1.RegisterNoteServiceServer(s, p.NoteServiceServer)
	gophkeeperv1.RegisterBinaryServiceServer(s, p.BinaryServiceServer)
	gophkeeperv1.RegisterCardServiceServer(s, p.CardServiceServer)
	reflection.Register(s)

	srv := &Server{
		s:      s,
		port:   p.Config.GRPC.Port,
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
