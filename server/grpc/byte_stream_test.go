package grpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/mock"
	"github.com/stretchr/testify/require"
	bs "google.golang.org/genproto/googleapis/bytestream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type writeTestStream struct {
	grpc.ServerStream
	ctx      context.Context
	requests []*bs.WriteRequest
	index    int
	response *bs.WriteResponse
	recvErr  error
	sendErr  error
}

func (s *writeTestStream) Context() context.Context { return s.ctx }
func (s *writeTestStream) Recv() (*bs.WriteRequest, error) {
	if s.index == len(s.requests) {
		if s.recvErr != nil {
			return nil, s.recvErr
		}
		return nil, io.EOF
	}
	req := s.requests[s.index]
	s.index++
	return req, nil
}
func (s *writeTestStream) SendAndClose(out *bs.WriteResponse) error {
	s.response = out
	return s.sendErr
}

type readTestStream struct {
	grpc.ServerStream
	ctx       context.Context
	responses []*bs.ReadResponse
	sendErr   error
}

func (s *readTestStream) Context() context.Context { return s.ctx }
func (s *readTestStream) Send(out *bs.ReadResponse) error {
	if s.sendErr != nil {
		return s.sendErr
	}
	s.responses = append(s.responses, out)
	return nil
}

type readSeekTestFile struct {
	*bytes.Reader
	closed bool
}

func (f *readSeekTestFile) Close() error { f.closed = true; return nil }

func TestByteStreamWrite(t *testing.T) {
	for _, finish := range []bool{false, true} {
		t.Run(map[bool]string{false: "unfinished", true: "finished"}[finish], func(t *testing.T) {
			service := &mock.BinaryServiceMock{
				WriteChunkFunc: func(ctx context.Context, id, offset int64, data []byte, finished bool) (int64, error) {
					require.Equal(t, t.Context(), ctx)
					require.EqualValues(t, 7, id)
					return offset + int64(len(data)), nil
				},
			}
			stream := &writeTestStream{ctx: t.Context(), requests: []*bs.WriteRequest{
				{ResourceName: "binary/7", WriteOffset: 4, Data: []byte("hello")},
				{WriteOffset: 9, Data: []byte("world"), FinishWrite: finish},
			}}
			require.NoError(t, NewByteStreamServer(service).Write(stream))
			require.NotNil(t, stream.response)
			require.EqualValues(t, 14, stream.response.GetCommittedSize())
			calls := service.WriteChunkCalls()
			require.Len(t, calls, 2)
			require.EqualValues(t, 4, calls[0].Offset)
			require.Equal(t, []byte("hello"), calls[0].Data)
			require.False(t, calls[0].Finish)
			require.EqualValues(t, 9, calls[1].Offset)
			require.Equal(t, []byte("world"), calls[1].Data)
			require.Equal(t, finish, calls[1].Finish)
		})
	}
	t.Run("empty_final_chunk", func(t *testing.T) {
		service := &mock.BinaryServiceMock{WriteChunkFunc: func(_ context.Context, _ int64, offset int64, data []byte, finish bool) (int64, error) {
			require.Empty(t, data)
			require.True(t, finish)
			return offset, nil
		}}
		stream := &writeTestStream{ctx: t.Context(), requests: []*bs.WriteRequest{
			{ResourceName: "binary/7", WriteOffset: 10, FinishWrite: true},
		}}
		require.NoError(t, NewByteStreamServer(service).Write(stream))
		require.EqualValues(t, 10, stream.response.GetCommittedSize())
	})
}

func TestByteStreamWriteInvalidRequest(t *testing.T) {
	cases := []struct {
		name     string
		requests []*bs.WriteRequest
		calls    int
	}{
		{name: "empty_stream"},
		{name: "missing_resource", requests: []*bs.WriteRequest{{}}},
		{name: "wrong_resource", requests: []*bs.WriteRequest{{ResourceName: "note/7"}}},
		{name: "invalid_id", requests: []*bs.WriteRequest{{ResourceName: "binary/bad"}}},
		{name: "negative_offset", requests: []*bs.WriteRequest{{ResourceName: "binary/7", WriteOffset: -1}}},
		{name: "changed_resource", requests: []*bs.WriteRequest{
			{ResourceName: "binary/7"}, {ResourceName: "binary/8"},
		}, calls: 1},
		{name: "after_finish", requests: []*bs.WriteRequest{
			{ResourceName: "binary/7", FinishWrite: true}, {},
		}, calls: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mock.BinaryServiceMock{}
			stream := &writeTestStream{ctx: t.Context(), requests: tc.requests}
			requireGrpcError(t, NewByteStreamServer(service).Write(stream), codes.InvalidArgument)
			require.Len(t, service.WriteChunkCalls(), tc.calls)
			require.Nil(t, stream.response)
		})
	}
}

func TestByteStreamWriteErrors(t *testing.T) {
	for _, tc := range contentErrorCases() {
		t.Run(tc.name, func(t *testing.T) {
			service := &mock.BinaryServiceMock{WriteChunkFunc: func(context.Context, int64, int64, []byte, bool) (int64, error) {
				return 0, tc.err
			}}
			stream := &writeTestStream{ctx: t.Context(), requests: []*bs.WriteRequest{{ResourceName: "binary/7"}}}
			requireGrpcError(t, NewByteStreamServer(service).Write(stream), tc.code)
			require.Nil(t, stream.response)
		})
	}
	t.Run("receive_error", func(t *testing.T) {
		service := &mock.BinaryServiceMock{}
		stream := &writeTestStream{ctx: t.Context(), recvErr: status.Error(codes.Unavailable, "connection lost")}
		requireGrpcError(t, NewByteStreamServer(service).Write(stream), codes.Unavailable)
		require.Empty(t, service.WriteChunkCalls())
	})
	t.Run("interrupted_after_chunk", func(t *testing.T) {
		service := &mock.BinaryServiceMock{}
		stream := &writeTestStream{ctx: t.Context(), requests: []*bs.WriteRequest{{ResourceName: "binary/7"}},
			recvErr: status.Error(codes.Unavailable, "connection lost")}
		requireGrpcError(t, NewByteStreamServer(service).Write(stream), codes.Unavailable)
		require.Len(t, service.WriteChunkCalls(), 1)
		require.Nil(t, stream.response)
	})
	t.Run("send_error", func(t *testing.T) {
		stream := &writeTestStream{ctx: t.Context(), requests: []*bs.WriteRequest{{ResourceName: "binary/7"}},
			sendErr: status.Error(codes.Unavailable, "connection lost")}
		requireGrpcError(t, NewByteStreamServer(&mock.BinaryServiceMock{}).Write(stream), codes.Unavailable)
	})
}

func TestByteStreamRead(t *testing.T) {
	content := bytes.Repeat([]byte("abcdef"), 25000)
	cases := []struct {
		name          string
		offset, limit int64
	}{
		{"whole_file", 0, 0}, {"offset", 17, 0}, {"limited", 17, 123},
		{"limit_beyond_end", int64(len(content)) - 3, 100}, {"at_end", int64(len(content)), 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			file := &readSeekTestFile{Reader: bytes.NewReader(content)}
			service := &mock.BinaryServiceMock{OpenFunc: func(ctx context.Context, id int64) (io.ReadSeekCloser, int64, error) {
				require.Equal(t, t.Context(), ctx)
				require.EqualValues(t, 7, id)
				return file, int64(len(content)), nil
			}}
			stream := &readTestStream{ctx: t.Context()}
			err := NewByteStreamServer(service).Read(&bs.ReadRequest{ResourceName: "binary/7", ReadOffset: tc.offset, ReadLimit: tc.limit}, stream)
			require.NoError(t, err)
			require.True(t, file.closed)
			var got []byte
			for _, resp := range stream.responses {
				require.LessOrEqual(t, len(resp.GetData()), 64*1024)
				got = append(got, resp.GetData()...)
			}
			end := int64(len(content))
			if tc.limit > 0 && tc.limit < end-tc.offset {
				end = tc.offset + tc.limit
			}
			require.Equal(t, string(content[tc.offset:end]), string(got))
		})
	}
}

func TestByteStreamReadErrors(t *testing.T) {
	for _, tc := range contentErrorCases() {
		t.Run(tc.name, func(t *testing.T) {
			service := &mock.BinaryServiceMock{OpenFunc: func(context.Context, int64) (io.ReadSeekCloser, int64, error) {
				return nil, 0, tc.err
			}}
			requireGrpcError(t, NewByteStreamServer(service).Read(
				&bs.ReadRequest{ResourceName: "binary/7"}, &readTestStream{ctx: t.Context()}), tc.code)
		})
	}
	for _, tc := range []struct {
		name string
		req  *bs.ReadRequest
		code codes.Code
	}{
		{"invalid_resource", &bs.ReadRequest{ResourceName: "note/7"}, codes.InvalidArgument},
		{"invalid_id", &bs.ReadRequest{ResourceName: "binary/bad"}, codes.InvalidArgument},
		{"negative_offset", &bs.ReadRequest{ResourceName: "binary/7", ReadOffset: -1}, codes.OutOfRange},
		{"negative_limit", &bs.ReadRequest{ResourceName: "binary/7", ReadLimit: -1}, codes.OutOfRange},
	} {
		t.Run(tc.name, func(t *testing.T) {
			service := &mock.BinaryServiceMock{}
			requireGrpcError(t, NewByteStreamServer(service).Read(tc.req, &readTestStream{ctx: t.Context()}), tc.code)
			require.Empty(t, service.OpenCalls())
		})
	}
	for _, scenario := range []string{"offset_beyond_end", "truncated_file", "send_error"} {
		t.Run(scenario, func(t *testing.T) {
			file := &readSeekTestFile{Reader: bytes.NewReader([]byte("abc"))}
			size := int64(3)
			if scenario == "truncated_file" {
				size = 10
			}
			service := &mock.BinaryServiceMock{OpenFunc: func(context.Context, int64) (io.ReadSeekCloser, int64, error) { return file, size, nil }}
			stream := &readTestStream{ctx: t.Context()}
			req := &bs.ReadRequest{ResourceName: "binary/7"}
			code := codes.Internal
			if scenario == "offset_beyond_end" {
				req.ReadOffset = 4
				code = codes.OutOfRange
			}
			if scenario == "send_error" {
				stream.sendErr = status.Error(codes.Unavailable, "connection lost")
				code = codes.Unavailable
			}
			requireGrpcError(t, NewByteStreamServer(service).Read(req, stream), code)
			require.True(t, file.closed)
		})
	}
}

func TestByteStreamQueryWriteStatus(t *testing.T) {
	for _, complete := range []bool{false, true} {
		t.Run(map[bool]string{false: "pending", true: "complete"}[complete], func(t *testing.T) {
			service := &mock.BinaryServiceMock{WriteStatusFunc: func(ctx context.Context, id int64) (int64, bool, error) {
				require.Equal(t, t.Context(), ctx)
				require.EqualValues(t, 7, id)
				return 123, complete, nil
			}}
			out, err := NewByteStreamServer(service).QueryWriteStatus(t.Context(), &bs.QueryWriteStatusRequest{ResourceName: "binary/7"})
			require.NoError(t, err)
			require.EqualValues(t, 123, out.GetCommittedSize())
			require.Equal(t, complete, out.GetComplete())
		})
	}
	for _, tc := range contentErrorCases() {
		t.Run(tc.name, func(t *testing.T) {
			service := &mock.BinaryServiceMock{WriteStatusFunc: func(context.Context, int64) (int64, bool, error) { return 0, false, tc.err }}
			out, err := NewByteStreamServer(service).QueryWriteStatus(t.Context(), &bs.QueryWriteStatusRequest{ResourceName: "binary/7"})
			requireGrpcError(t, err, tc.code)
			require.Nil(t, out)
		})
	}
	t.Run("invalid_resource", func(t *testing.T) {
		service := &mock.BinaryServiceMock{}
		_, err := NewByteStreamServer(service).QueryWriteStatus(t.Context(), &bs.QueryWriteStatusRequest{ResourceName: "wrong"})
		requireGrpcError(t, err, codes.InvalidArgument)
		require.Empty(t, service.WriteStatusCalls())
	})
}

func contentErrorCases() []struct {
	name string
	err  error
	code codes.Code
} {
	return []struct {
		name string
		err  error
		code codes.Code
	}{
		{"not_found", server.ErrDataNotFound, codes.NotFound},
		{"incomplete", server.ErrUploadIncomplete, codes.FailedPrecondition},
		{"service_error", errors.New("storage error"), codes.Internal},
		{"canceled", context.Canceled, codes.Canceled},
		{"deadline", context.DeadlineExceeded, codes.DeadlineExceeded},
	}
}
