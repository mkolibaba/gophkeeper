package interceptors

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
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

		baOpt, err := newBearerAccessOption(userService)
		if err != nil {
			return fmt.Errorf("auth interceptor: %w", err)
		}
		opts = append(opts, baOpt)

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

		baOpt, err := newBearerAccessOption(userService)
		if err != nil {
			return nil, fmt.Errorf("auth interceptor: %w", err)
		}
		opts = append(opts, baOpt)

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func newBearerAccessOption(userService client.UserService) (grpc.CallOption, error) {
	token := userService.GetBearerToken()
	if token == nil {
		return nil, fmt.Errorf("bearer token not found in session")
	}
	return grpc.PerRPCCredentials(bearerAccess{*token}), nil
}

type bearerAccess struct {
	token string
}

func (b bearerAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + b.token,
	}, nil
}

func (b bearerAccess) RequireTransportSecurity() bool {
	return false
}
