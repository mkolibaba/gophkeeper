package grpc

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/client/grpc/mock"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"testing"
)

func TestAuthorizationAuthorize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := NewAuthorizationService(
			&mock.AuthorizationServiceClientMock{
				AuthorizeFunc: func(ctx context.Context, in *gophkeeperv1.UserCredentials, opts ...grpc.CallOption) (*gophkeeperv1.TokenResponse, error) {
					var out gophkeeperv1.TokenResponse
					out.SetToken("cool token")
					return &out, nil
				},
			},
			log.New(io.Discard),
		)

		resp, err := srv.Authorize(t.Context(), "testuser", "123")
		require.NoError(t, err)
		require.NotEmpty(t, resp)
	})
	t.Run("fail", func(t *testing.T) {
		srv := NewAuthorizationService(
			&mock.AuthorizationServiceClientMock{
				AuthorizeFunc: func(ctx context.Context, in *gophkeeperv1.UserCredentials, opts ...grpc.CallOption) (*gophkeeperv1.TokenResponse, error) {
					return nil, status.Error(codes.Internal, "some error")
				},
			},
			log.New(io.Discard),
		)

		_, err := srv.Authorize(t.Context(), "testuser", "123")
		require.Error(t, err)
	})
}
