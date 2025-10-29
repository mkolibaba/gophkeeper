package grpc

import (
	"github.com/go-playground/validator/v10"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func requireGrpcError(t *testing.T, err error, code codes.Code) {
	t.Helper()

	require.Error(t, err)
	s, ok := status.FromError(err)
	require.True(t, ok, "error should be a grpc Status")
	require.Equal(t, s.Code(), code)
}

func newTestValidator(t *testing.T) *validator.Validate {
	t.Helper()
	v := validator.New()
	err := server.RegisterDataValidationRules(v)
	require.NoError(t, err)
	RegisterValidationRules(v)
	return v
}
