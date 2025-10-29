package interceptors

import (
	"bytes"
	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/interop/grpc_testing"
	"io"
	"testing"
)

func TestLoggedAttributesUnary(t *testing.T) {
	var out bytes.Buffer
	i := NewLoggerInterceptor(log.New(&out))

	stubServer := newStubServer()
	stubServer.startServer(grpc.UnaryInterceptor(i.Unary))

	if err := stubServer.startClient(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(stubServer.stop)

	_, err := stubServer.client.UnaryCall(t.Context(), &grpc_testing.SimpleRequest{})
	require.NoError(t, err)

	callLog := out.String()
	require.NotEmpty(t, callLog)
	// Проверяем несколько полей.
	require.Regexp(t, `status_code`, callLog)
	require.Regexp(t, `method=\S+UnaryCall`, callLog)
}

func TestLoggedAttributesStream(t *testing.T) {
	var out bytes.Buffer
	i := NewLoggerInterceptor(log.New(&out))

	stubServer := newStubServer()
	stubServer.startServer(grpc.StreamInterceptor(i.Stream))

	if err := stubServer.startClient(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(stubServer.stop)

	resp, err := stubServer.client.StreamingOutputCall(t.Context(), &grpc_testing.StreamingOutputCallRequest{})
	require.NoError(t, err)

	_, err = resp.Recv() // читаем сообщение
	require.NoError(t, err)
	_, err = resp.Recv() // читаем io.EOF
	require.Equal(t, io.EOF, err)

	callLog := out.String()
	require.NotEmpty(t, callLog)
	// Проверяем несколько полей.
	require.Regexp(t, `status_code`, callLog)
	require.Regexp(t, `method=\S+StreamingOutputCall`, callLog)
}
