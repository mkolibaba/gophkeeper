package interceptors

import (
	"context"
	"encoding/base64"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/grpc/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"slices"
	"strings"
)

var skip = []string{
	gophkeeperv1.AuthorizationService_Authorize_FullMethodName,
	gophkeeperv1.AuthorizationService_Register_FullMethodName,
	grpc_reflection_v1.ServerReflection_ServerReflectionInfo_FullMethodName,
	grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfo_FullMethodName,
}

func UnaryAuth(authService *server.AuthService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if slices.Contains(skip, info.FullMethod) {
			return handler(ctx, req)
		}

		authorization := metadata.ValueFromIncomingContext(ctx, "authorization")
		if len(authorization) == 0 {
			return nil, status.Error(codes.Unauthenticated, "credentials are not provided")
		}

		credentials := authorization[0]
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(credentials, "Basic "))
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
		}

		split := strings.Split(string(decoded), ":")
		if len(split) != 2 {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials format")
		}

		login := split[0]
		password := split[1]

		if err := authService.Authorize(ctx, login, password); err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		userCtx := utils.NewContextWithUser(ctx, login)
		return handler(userCtx, req)
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func StreamAuth(authService *server.AuthService) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if slices.Contains(skip, info.FullMethod) {
			return handler(srv, ss)
		}

		authorization := metadata.ValueFromIncomingContext(ss.Context(), "authorization")
		if len(authorization) == 0 {
			return status.Errorf(codes.Unauthenticated, "credentials are not provided")
		}

		credentials := authorization[0]
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(credentials, "Basic "))
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
		}

		split := strings.Split(string(decoded), ":")
		if len(split) != 2 {
			return status.Errorf(codes.Unauthenticated, "invalid credentials format")
		}

		login := split[0]
		password := split[1]

		if err := authService.Authorize(ss.Context(), login, password); err != nil {
			return status.Error(codes.Unauthenticated, err.Error())
		}

		userCtx := utils.NewContextWithUser(ss.Context(), login)
		wss := &wrappedServerStream{ServerStream: ss, ctx: userCtx}

		return handler(srv, wss)
	}
}
