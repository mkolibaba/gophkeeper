package grpc

import (
	"github.com/mkolibaba/gophkeeper/client/core"
	"github.com/mkolibaba/gophkeeper/client/core/grpc/interceptors"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionParams struct {
	fx.In

	Config      *core.Config
	UserService core.UserService
}

func NewConnection(p ConnectionParams) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		p.Config.GRPC.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryAuth(p.UserService)),
		grpc.WithStreamInterceptor(interceptors.StreamAuth(p.UserService)),
	)
	return conn, err
}
