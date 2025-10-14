package grpc

import (
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/client/grpc/interceptors"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionParams struct {
	fx.In

	Config      *Config
	UserService client.UserService
}

func NewConnection(p ConnectionParams) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		p.Config.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryAuth(p.UserService)),
		grpc.WithStreamInterceptor(interceptors.StreamAuth(p.UserService)),
	)
	return conn, err
}
