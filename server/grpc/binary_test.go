package grpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// mockUploadStream for testing client-side streaming
type mockUploadStream struct {
	grpc.ClientStream
	ctx      context.Context
	requests []*gophkeeperv1.SaveBinaryRequest
	recvIdx  int
}

func (s *mockUploadStream) Context() context.Context { return s.ctx }
func (s *mockUploadStream) Recv() (*gophkeeperv1.SaveBinaryRequest, error) {
	if s.recvIdx >= len(s.requests) {
		return nil, io.EOF
	}
	req := s.requests[s.recvIdx]
	s.recvIdx++
	return req, nil
}
func (s *mockUploadStream) SendAndClose(_ *empty.Empty) error { return nil }
func (s *mockUploadStream) SendHeader(_ metadata.MD) error    { return nil }
func (s *mockUploadStream) SetHeader(_ metadata.MD) error     { return nil }
func (s *mockUploadStream) SetTrailer(_ metadata.MD)          {}
func (s *mockUploadStream) CloseSend() error                  { return nil }
func (s *mockUploadStream) Header() (metadata.MD, error)      { return nil, nil }
func (s *mockUploadStream) Trailer() metadata.MD              { return nil }

// mockDownloadStream for testing server-side streaming
type mockDownloadStream struct {
	grpc.ServerStream
	ctx       context.Context
	responses []*gophkeeperv1.DownloadBinaryResponse
}

func (s *mockDownloadStream) Context() context.Context { return s.ctx }
func (s *mockDownloadStream) Send(resp *gophkeeperv1.DownloadBinaryResponse) error {
	s.responses = append(s.responses, resp)
	return nil
}
func (s *mockDownloadStream) SendHeader(_ metadata.MD) error { return nil }
func (s *mockDownloadStream) SetHeader(_ metadata.MD) error  { return nil }
func (s *mockDownloadStream) SetTrailer(_ metadata.MD)       {}

func TestBinaryUpload(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mock.BinaryServiceMock{
			CreateFunc: func(ctx context.Context, data server.ReadableBinaryData) error {
				return nil
			},
		}
		srv := createBinaryServiceServer(t, service)

		var req1 gophkeeperv1.SaveBinaryRequest
		req1.SetName("testfile")
		req1.SetFilename("test.txt")
		req1.SetSize(10)
		var chunk1 gophkeeperv1.FileChunk
		chunk1.SetData([]byte("hello"))
		req1.SetChunk(&chunk1)

		var req2 gophkeeperv1.SaveBinaryRequest
		var chunk2 gophkeeperv1.FileChunk
		chunk2.SetData([]byte("world"))
		req2.SetChunk(&chunk2)

		stream := &mockUploadStream{
			ctx:      t.Context(),
			requests: []*gophkeeperv1.SaveBinaryRequest{&req1, &req2},
		}

		err := srv.Upload(stream)
		require.NoError(t, err)
	})

	t.Run("validation_error", func(t *testing.T) {
		srv := createBinaryServiceServer(t, &mock.BinaryServiceMock{})
		var req gophkeeperv1.SaveBinaryRequest
		var chunk gophkeeperv1.FileChunk
		chunk.SetData([]byte("hello"))
		req.SetChunk(&chunk)

		stream := &mockUploadStream{
			ctx:      t.Context(),
			requests: []*gophkeeperv1.SaveBinaryRequest{&req},
		}
		err := srv.Upload(stream)
		requireGrpcError(t, err, codes.InvalidArgument)
	})

	t.Run("service_error", func(t *testing.T) {
		service := &mock.BinaryServiceMock{
			CreateFunc: func(ctx context.Context, data server.ReadableBinaryData) error {
				return fmt.Errorf("db error")
			},
		}
		srv := createBinaryServiceServer(t, service)

		var req gophkeeperv1.SaveBinaryRequest
		req.SetName("testfile")
		req.SetFilename("test.txt")
		req.SetSize(5)
		var chunk gophkeeperv1.FileChunk
		chunk.SetData([]byte("hello"))
		req.SetChunk(&chunk)

		stream := &mockUploadStream{
			ctx:      t.Context(),
			requests: []*gophkeeperv1.SaveBinaryRequest{&req},
		}
		err := srv.Upload(stream)
		requireGrpcError(t, err, codes.Internal)
	})
}

func TestBinaryDownload(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fileContent := "this is the file content"
		service := &mock.BinaryServiceMock{
			GetFunc: func(ctx context.Context, id int64) (*server.ReadableBinaryData, error) {
				return &server.ReadableBinaryData{
					BinaryData: server.BinaryData{
						Name:     "testfile",
						Filename: "test.txt",
						Size:     int64(len(fileContent)),
					},
					DataReader: io.NopCloser(bytes.NewReader([]byte(fileContent))),
				}, nil
			},
		}
		srv := createBinaryServiceServer(t, service)

		var req gophkeeperv1.DownloadBinaryRequest
		req.SetId(1)

		stream := &mockDownloadStream{ctx: t.Context()}
		err := srv.Download(&req, stream)
		require.NoError(t, err)
		require.NotEmpty(t, stream.responses)

		var downloadedContent []byte
		for _, resp := range stream.responses {
			downloadedContent = append(downloadedContent, resp.GetChunk().GetData()...)
		}
		require.Equal(t, fileContent, string(downloadedContent))
	})

	t.Run("not_found", func(t *testing.T) {
		service := &mock.BinaryServiceMock{
			GetFunc: func(ctx context.Context, id int64) (*server.ReadableBinaryData, error) {
				return nil, server.ErrDataNotFound
			},
		}
		srv := createBinaryServiceServer(t, service)
		var req gophkeeperv1.DownloadBinaryRequest
		req.SetId(1)
		stream := &mockDownloadStream{ctx: t.Context()}
		err := srv.Download(&req, stream)
		requireGrpcError(t, err, codes.NotFound)
	})

	t.Run("service_error", func(t *testing.T) {
		service := &mock.BinaryServiceMock{
			GetFunc: func(ctx context.Context, id int64) (*server.ReadableBinaryData, error) {
				return nil, fmt.Errorf("some error")
			},
		}
		srv := createBinaryServiceServer(t, service)
		var req gophkeeperv1.DownloadBinaryRequest
		req.SetId(1)
		stream := &mockDownloadStream{ctx: t.Context()}
		err := srv.Download(&req, stream)
		requireGrpcError(t, err, codes.Internal)
	})
}

func TestBinaryUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createBinaryServiceServer(t, &mock.BinaryServiceMock{})

		var in gophkeeperv1.UpdateBinaryRequest
		in.SetId(1)
		in.SetName("new binary name")

		_, err := srv.Update(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createBinaryServiceServer(t, &mock.BinaryServiceMock{})

		var in gophkeeperv1.UpdateBinaryRequest
		in.SetName("new binary name")

		_, err := srv.Update(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
}

func TestBinaryRemove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createBinaryServiceServer(t, &mock.BinaryServiceMock{})

		var in gophkeeperv1.RemoveDataRequest
		in.SetId(1)

		_, err := srv.Remove(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createBinaryServiceServer(t, &mock.BinaryServiceMock{})

		var in gophkeeperv1.RemoveDataRequest

		_, err := srv.Remove(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
	t.Run("not_found", func(t *testing.T) {
		service := &mock.BinaryServiceMock{
			RemoveFunc: func(ctx context.Context, id int64) error {
				return server.ErrDataNotFound
			},
		}
		srv := createBinaryServiceServer(t, service)

		var in gophkeeperv1.RemoveDataRequest
		in.SetId(1)

		_, err := srv.Remove(t.Context(), &in)
		requireGrpcError(t, err, codes.NotFound)
	})
}

func TestBinaryGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mock.BinaryServiceMock{
			GetAllFunc: func(ctx context.Context) ([]server.BinaryData, error) {
				return []server.BinaryData{
					{ID: 1, Name: "binary1"},
					{ID: 2, Name: "binary2"},
				}, nil
			},
		}
		srv := createBinaryServiceServer(t, service)
		resp, err := srv.GetAll(t.Context(), nil)
		require.NoError(t, err)
		require.Len(t, resp.GetResult(), 2)
	})
	t.Run("db_error", func(t *testing.T) {
		service := &mock.BinaryServiceMock{
			GetAllFunc: func(ctx context.Context) ([]server.BinaryData, error) {
				return nil, fmt.Errorf("db error")
			},
		}
		srv := createBinaryServiceServer(t, service)
		_, err := srv.GetAll(t.Context(), nil)
		requireGrpcError(t, err, codes.Internal)
	})
}

func createBinaryServiceServer(t *testing.T, binaryService server.BinaryService) *BinaryServiceServer {
	return NewBinaryServiceServer(binaryService, newTestValidator(t), log.New(io.Discard))
}
