package interceptors

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

type AuthInterceptor struct {
	config *server.Config
}

func NewAuthInterceptor(config *server.Config) *AuthInterceptor {
	return &AuthInterceptor{config: config}
}

func (i *AuthInterceptor) Unary(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if slices.Contains(skip, info.FullMethod) {
		return handler(ctx, req)
	}

	sub, err := i.getSub(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	userCtx := utils.NewContextWithUser(ctx, sub)
	return handler(userCtx, req)
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func (i *AuthInterceptor) Stream(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if slices.Contains(skip, info.FullMethod) {
		return handler(srv, ss)
	}

	sub, err := i.getSub(ss.Context())
	if err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}

	userCtx := utils.NewContextWithUser(ss.Context(), sub)
	wss := &wrappedServerStream{ServerStream: ss, ctx: userCtx}

	return handler(srv, wss)
}

func (i *AuthInterceptor) getSub(ctx context.Context) (string, error) {
	authorization := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(authorization) == 0 {
		return "", fmt.Errorf("credentials are not provided")
	}

	tokenString := strings.TrimPrefix(authorization[0], "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing method")
		}
		return i.config.GetJWTSecret(), nil
	})
	if err != nil {
		return "", fmt.Errorf("bearer token: %w", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("bearer token: invalid")
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("jwt claims: %w", err)
	}

	return sub, nil
}
