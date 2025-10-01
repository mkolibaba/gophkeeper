package interceptors

import (
	"context"
	"encoding/base64"
	"github.com/mkolibaba/gophkeeper/internal/server/grpc/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func UnaryAuth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		authorization := metadata.ValueFromIncomingContext(ctx, "authorization")
		if len(authorization) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "credentials are not provided")
		}

		credentials := authorization[0]
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(credentials, "Basic "))
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
		}

		split := strings.Split(string(decoded), ":")
		if len(split) != 2 {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials format")
		}

		login := split[0]
		password := split[1]

		if login != "demo" || password != "demo" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid login/password")
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

func StreamAuth() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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

		if login != "demo" || password != "demo" {
			return status.Errorf(codes.Unauthenticated, "invalid login/password")
		}

		userCtx := utils.NewContextWithUser(ss.Context(), login)
		wss := &wrappedServerStream{ServerStream: ss, ctx: userCtx}

		return handler(srv, wss)
	}
}
