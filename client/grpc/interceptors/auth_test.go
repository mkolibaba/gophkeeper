package interceptors

import (
	"context"
	"github.com/mkolibaba/gophkeeper/client/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/interop/grpc_testing"
	"testing"
)

func TestAuthUnary(t *testing.T) {
	stubServer := newStubServer()
	stubServer.UnaryCallF = func(ctx context.Context, request *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) {
		return &grpc_testing.SimpleResponse{
			Payload: &grpc_testing.Payload{
				Body: []byte("hi!"),
			},
		}, nil
	}
	stubServer.startServer()

	userServiceMock := &mock.UserServiceMock{}
	if err := stubServer.startClient(grpc.WithUnaryInterceptor(UnaryAuth(userServiceMock))); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(stubServer.stop)

	t.Run("success", func(t *testing.T) {
		userServiceMock.GetBearerTokenFunc = func() string {
			return "cool token"
		}

		resp, err := stubServer.client.UnaryCall(t.Context(), &grpc_testing.SimpleRequest{})
		require.NoError(t, err)

		require.Equal(t, "hi!", string(resp.Payload.Body))
	})
	t.Run("invalid_token", func(t *testing.T) {
		userServiceMock.GetBearerTokenFunc = func() string {
			return ""
		}

		_, err := stubServer.client.UnaryCall(t.Context(), &grpc_testing.SimpleRequest{})
		require.Error(t, err)
	})
}

func TestAuthStream(t *testing.T) {
	stubServer := newStubServer()
	stubServer.startServer()

	userServiceMock := &mock.UserServiceMock{
		GetBearerTokenFunc: func() string {
			return "cool token"
		},
	}
	if err := stubServer.startClient(grpc.WithStreamInterceptor(StreamAuth(userServiceMock))); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(stubServer.stop)

	_, err := stubServer.client.StreamingOutputCall(t.Context(), &grpc_testing.StreamingOutputCallRequest{})
	require.NoError(t, err)
}
