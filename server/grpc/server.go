package grpc

import (
	"fmt"
	"net"

	"charm.land/log/v2"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/grpc/interceptors"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	s      *grpc.Server
	port   string
	logger *log.Logger
}

type ServerParams struct {
	fx.In

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

	return &Server{
		s:      s,
		port:   p.Config.GRPC.Port,
		logger: p.Logger,
	}
}

func (s *Server) Listen() (net.Listener, error) {
	return net.Listen("tcp", ":"+s.port)
}

func (s *Server) Start(ln net.Listener) error {
	s.logger.Info("running grpc server", "port", s.port)

	if err := s.s.Serve(ln); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	s.logger.Info("server stopped")

	return nil
}

func (s *Server) Stop() {
	s.s.GracefulStop()
}
