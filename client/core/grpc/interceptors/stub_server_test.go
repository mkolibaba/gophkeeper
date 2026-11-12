package interceptors

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"google.golang.org/grpc/test/bufconn"
	"net"
)

type stubServer struct {
	grpc_testing.TestServiceServer

	listener *bufconn.Listener
	resolver *manual.Resolver
	conn     *grpc.ClientConn
	client   grpc_testing.TestServiceClient

	UnaryCallF func(context.Context, *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error)
}

func (s *stubServer) UnaryCall(ctx context.Context, in *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) {
	return s.UnaryCallF(ctx, in)
}

func (*stubServer) StreamingOutputCall(
	_ *grpc_testing.StreamingOutputCallRequest,
	out grpc.ServerStreamingServer[grpc_testing.StreamingOutputCallResponse],
) error {
	return out.Send(&grpc_testing.StreamingOutputCallResponse{
		Payload: &grpc_testing.Payload{
			Body: []byte("Hello!"),
		},
	})
}

func newStubServer() *stubServer {
	listen := bufconn.Listen(20)

	r := manual.NewBuilderWithScheme("whatever")
	r.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: listen.Addr().String()}}})

	return &stubServer{
		listener: listen,
		resolver: r,
		UnaryCallF: func(context.Context, *grpc_testing.SimpleRequest) (*grpc_testing.SimpleResponse, error) {
			return &grpc_testing.SimpleResponse{
				Payload: &grpc_testing.Payload{
					Body: []byte("Hello!"),
				},
			}, nil
		},
	}
}

func (s *stubServer) startServer(opts ...grpc.ServerOption) {
	srv := grpc.NewServer(opts...)
	grpc_testing.RegisterTestServiceServer(srv, s)

	go func() {
		defer srv.Stop()
		srv.Serve(s.listener)
	}()
}

func (s *stubServer) startClient(opts ...grpc.DialOption) (err error) {
	os := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return s.listener.DialContext(ctx)
		}),
		grpc.WithResolvers(s.resolver),
	}
	for _, opt := range opts {
		os = append(os, opt)
	}

	s.conn, err = grpc.NewClient(
		fmt.Sprintf("%s:///%s", s.resolver.Scheme(), s.listener.Addr().String()),
		os...,
	)
	if err != nil {
		return
	}

	s.client = grpc_testing.NewTestServiceClient(s.conn)
	return
}

func (s *stubServer) stop() {
	s.listener.Close()
	s.conn.Close()
}
