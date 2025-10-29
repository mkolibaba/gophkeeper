package grpc

import (
	"context"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"testing"
)

func TestAuthorize(t *testing.T) {
	cases := map[string]struct {
		login    string
		password string
		checks   func(out *gophkeeperv1.TokenResponse, err error)
	}{
		"success": {
			login:    "alice",
			password: "123",
			checks: func(out *gophkeeperv1.TokenResponse, err error) {
				t.Helper()
				require.NoError(t, err)
				require.NotEmpty(t, out.GetToken())
			},
		},
		"invalid_login": {
			login:    "bob",
			password: "123",
			checks: func(out *gophkeeperv1.TokenResponse, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
		"invalid_password": {
			login:    "alice",
			password: "1234",
			checks: func(out *gophkeeperv1.TokenResponse, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
		"no_password": {
			login:    "alice",
			password: "",
			checks: func(out *gophkeeperv1.TokenResponse, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
	}

	userServiceMock := &mock.UserServiceMock{
		GetFunc: func(ctx context.Context, login string) (*server.User, error) {
			if login != "alice" {
				return nil, server.ErrUserNotFound
			}

			hash, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
			require.NoError(t, err)
			return &server.User{
				Login:    "alice",
				Password: string(hash),
			}, nil
		},
	}
	authorizationServiceMock := &mock.AuthorizationServiceMock{
		AuthorizeFunc: func(ctx context.Context, login string) (string, error) {
			return "some token", nil
		},
	}

	srv := createAuthorizationServiceServer(t, userServiceMock, authorizationServiceMock)

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var in gophkeeperv1.UserCredentials
			in.SetLogin(c.login)
			in.SetPassword(c.password)
			out, err := srv.Authorize(t.Context(), &in)
			c.checks(out, err)
		})
	}
}

func TestRegister(t *testing.T) {
	userServiceMock := &mock.UserServiceMock{
		SaveFunc: func(ctx context.Context, user server.User) error {
			if user.Login == "alice" {
				return server.ErrUserAlreadyExists
			}
			return nil
		},
	}
	authorizationServiceMock := &mock.AuthorizationServiceMock{
		AuthorizeFunc: func(ctx context.Context, login string) (string, error) {
			return "some token", nil
		},
	}

	srv := createAuthorizationServiceServer(t, userServiceMock, authorizationServiceMock)

	cases := map[string]struct {
		login    string
		password string
		checks   func(out *gophkeeperv1.TokenResponse, err error)
	}{
		"success": {
			login:    "bob",
			password: "123",
			checks: func(out *gophkeeperv1.TokenResponse, err error) {
				t.Helper()
				require.NoError(t, err)
			},
		},
		"user_exists": {
			login:    "alice",
			password: "123",
			checks: func(out *gophkeeperv1.TokenResponse, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrInvalidCredentials)
			},
		},
		"no_password": {
			login:    "bob",
			password: "",
			checks: func(out *gophkeeperv1.TokenResponse, err error) {
				t.Helper()
				requireGrpcError(t, err, codes.InvalidArgument)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var in gophkeeperv1.UserCredentials
			in.SetLogin(c.login)
			in.SetPassword(c.password)
			out, err := srv.Register(t.Context(), &in)
			c.checks(out, err)
		})
	}
}

func createAuthorizationServiceServer(
	t *testing.T,
	userService server.UserService,
	authorizationService server.AuthorizationService,
) *AuthorizationServiceServer {
	return NewAuthorizationServiceServer(userService, authorizationService, newTestValidator(t))
}
