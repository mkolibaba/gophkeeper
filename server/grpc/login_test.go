package grpc

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestLoginSave(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createLoginServiceServer(t, &mock.LoginServiceMock{})

		var in gophkeeperv1.Login
		in.SetName("new login")
		in.SetLogin("user")
		in.SetPassword("pass")

		_, err := srv.Save(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createLoginServiceServer(t, &mock.LoginServiceMock{})

		var in gophkeeperv1.Login

		_, err := srv.Save(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
	t.Run("db_error", func(t *testing.T) {
		loginServiceMock := &mock.LoginServiceMock{
			CreateFunc: func(_ context.Context, _ server.LoginData) error {
				return fmt.Errorf("some error")
			},
		}
		srv := createLoginServiceServer(t, loginServiceMock)

		var in gophkeeperv1.Login
		in.SetName("new login")
		in.SetLogin("user")
		in.SetPassword("pass")

		_, err := srv.Save(t.Context(), &in)
		requireGrpcError(t, err, codes.Internal)
	})
}

func TestLoginUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createLoginServiceServer(t, &mock.LoginServiceMock{})

		var in gophkeeperv1.Login
		in.SetId(1)
		in.SetName("new login name")

		_, err := srv.Update(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createLoginServiceServer(t, &mock.LoginServiceMock{})

		var in gophkeeperv1.Login
		in.SetName("new login name")

		_, err := srv.Update(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
}

func TestLoginRemove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createLoginServiceServer(t, &mock.LoginServiceMock{})

		var in gophkeeperv1.RemoveDataRequest
		in.SetId(1)

		_, err := srv.Remove(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createLoginServiceServer(t, &mock.LoginServiceMock{})

		var in gophkeeperv1.RemoveDataRequest

		_, err := srv.Remove(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
	t.Run("not_found", func(t *testing.T) {
		service := &mock.LoginServiceMock{
			RemoveFunc: func(ctx context.Context, id int64) error {
				return server.ErrDataNotFound
			},
		}
		srv := createLoginServiceServer(t, service)

		var in gophkeeperv1.RemoveDataRequest
		in.SetId(1)

		_, err := srv.Remove(t.Context(), &in)
		requireGrpcError(t, err, codes.NotFound)
	})
}

func TestLoginGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mock.LoginServiceMock{
			GetAllFunc: func(ctx context.Context) ([]server.LoginData, error) {
				return []server.LoginData{
					{ID: 1, Name: "login1"},
					{ID: 2, Name: "login2"},
				}, nil
			},
		}
		srv := createLoginServiceServer(t, service)
		resp, err := srv.GetAll(t.Context(), nil)
		require.NoError(t, err)
		require.Len(t, resp.GetResult(), 2)
	})
	t.Run("db_error", func(t *testing.T) {
		service := &mock.LoginServiceMock{
			GetAllFunc: func(ctx context.Context) ([]server.LoginData, error) {
				return nil, fmt.Errorf("db error")
			},
		}
		srv := createLoginServiceServer(t, service)
		_, err := srv.GetAll(t.Context(), nil)
		requireGrpcError(t, err, codes.Internal)
	})
}

func createLoginServiceServer(t *testing.T, loginService server.LoginService) *LoginServiceServer {
	return NewLoginServiceServer(loginService, newTestValidator(t), log.New(io.Discard))
}
