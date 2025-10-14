package interceptors

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"google.golang.org/grpc"
	"slices"
)

var skip = []string{
	gophkeeperv1.AuthorizationService_Authorize_FullMethodName,
	gophkeeperv1.AuthorizationService_Register_FullMethodName,
}

func UnaryAuth(userService client.UserService) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if slices.Contains(skip, method) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		user := userService.Get()
		if user == nil {
			return fmt.Errorf("user not found in session")
		}

		opts = append(opts, grpc.PerRPCCredentials(basicAccess{login: user.Login, password: user.Password}))

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func StreamAuth(userService client.UserService) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		if slices.Contains(skip, method) {
			return streamer(ctx, desc, cc, method, opts...)
		}

		user := userService.Get()
		if user == nil {
			return nil, fmt.Errorf("user not found in session")
		}

		opts = append(opts, grpc.PerRPCCredentials(basicAccess{login: user.Login, password: user.Password}))

		return streamer(ctx, desc, cc, method, opts...)
	}
}

type basicAccess struct {
	login    string
	password string
}

func (b basicAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(b.login+":"+b.password)),
	}, nil
}

func (b basicAccess) RequireTransportSecurity() bool {
	return false
}
