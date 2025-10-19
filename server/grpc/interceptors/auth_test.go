package interceptors

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"testing"
	"time"
)

const (
	// sub = "testuser", secret = "mysecret"
	validJWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ0MTE0OTQwMDAsImlzcyI6ImdvcGhrZWVwZXIiLCJzdWIiOiJ0ZXN0dXNlciJ9.dVWN19MmRr_zGihwS9dvedqOlAIv1zCHnRdC1TS8pGU"
)

func TestAuthUnary(t *testing.T) {
	config := &server.Config{
		JWT: struct {
			Secret string
			TTL    time.Duration
		}{Secret: "mysecret", TTL: 10 * time.Hour},
	}

	i := NewAuthInterceptor(config)

	stubServer := newStubServer()
	stubServer.UnaryCallF = func(ctx context.Context, request *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) {
		return &grpc_testing.SimpleResponse{
			Payload: &grpc_testing.Payload{
				Body: []byte(fmt.Sprintf("Hello, %s!", server.UserFromContext(ctx))),
			},
		}, nil
	}
	stubServer.startServer(grpc.UnaryInterceptor(i.Unary))

	if err := stubServer.startClient(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(stubServer.stop)

	t.Run("success", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(t.Context(), metadata.Pairs("authorization", "Bearer "+validJWT))

		resp, err := stubServer.client.UnaryCall(ctx, &grpc_testing.SimpleRequest{})
		require.NoError(t, err)

		require.Equal(t, "Hello, testuser!", string(resp.Payload.Body))
	})
	t.Run("no_auth", func(t *testing.T) {
		_, err := stubServer.client.UnaryCall(t.Context(), &grpc_testing.SimpleRequest{})
		require.Error(t, err)

		s, ok := status.FromError(err)
		require.True(t, ok, "error should be a grpc Status")
		require.Equal(t, s.Code(), codes.Unauthenticated)
	})
	t.Run("invalid_token", func(t *testing.T) {
		ctx := metadata.NewOutgoingContext(t.Context(), metadata.Pairs("authorization", "heh"))

		_, err := stubServer.client.UnaryCall(ctx, &grpc_testing.SimpleRequest{})
		require.Error(t, err)

		s, ok := status.FromError(err)
		require.True(t, ok, "error should be a grpc Status")
		require.Equal(t, s.Code(), codes.Unauthenticated)
	})
}

func TestAuthStream(t *testing.T) {
	config := &server.Config{
		JWT: struct {
			Secret string
			TTL    time.Duration
		}{Secret: "mysecret", TTL: 10 * time.Hour},
	}

	i := NewAuthInterceptor(config)

	stubServer := newStubServer()
	stubServer.startServer(grpc.StreamInterceptor(i.Stream))

	if err := stubServer.startClient(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(stubServer.stop)

	ctx := metadata.NewOutgoingContext(t.Context(), metadata.Pairs("authorization", validJWT))

	resp, err := stubServer.client.StreamingOutputCall(ctx, &grpc_testing.StreamingOutputCallRequest{})
	require.NoError(t, err)

	_, err = resp.Recv() // читаем сообщение
	require.NoError(t, err)
	_, err = resp.Recv() // читаем io.EOF
	require.Equal(t, io.EOF, err)
}
