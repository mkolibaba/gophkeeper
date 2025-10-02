package grpc

import (
	"context"
	"encoding/base64"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type basicAccess struct {
}

func (b basicAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	login := "demo"
	password := "demo"

	authorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password))

	return map[string]string{
		"authorization": authorization,
	}, nil
}

func (b basicAccess) RequireTransportSecurity() bool {
	return false // TODO
}

func NewConnection(cfg *Config) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		cfg.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// TODO: https://github.com/grpc/grpc-go/blob/master/examples/features/authz/client/main.go
		//  можно устанавливать токен прямо во время запроса
		grpc.WithPerRPCCredentials(basicAccess{}),
	)
	return conn, err
}
