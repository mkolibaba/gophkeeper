package grpc

import (
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/grpc/interceptors"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionParams struct {
	fx.In

	Config  *Config
	Session *client.Session
}

func NewConnection(p ConnectionParams) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		p.Config.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryAuth(p.Session)),
		grpc.WithStreamInterceptor(interceptors.StreamAuth(p.Session)),
	)
	return conn, err
}
