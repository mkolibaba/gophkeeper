package grpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mkolibaba/gophkeeper/server"
	"google.golang.org/genproto/googleapis/bytestream"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ByteStreamServer struct {
	bytestream.UnimplementedByteStreamServer
	binaryService server.BinaryService
}

func NewByteStreamServer(binaryService server.BinaryService) *ByteStreamServer {
	return &ByteStreamServer{
		binaryService: binaryService,
	}
}

func (s *ByteStreamServer) Read(req *bytestream.ReadRequest, stream bytestream.ByteStream_ReadServer) error {
	id, err := s.parseID(req.GetResourceName())
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	offset, limit := req.GetReadOffset(), req.GetReadLimit()
	if offset < 0 {
		return status.Error(codes.OutOfRange, "negative read offset")
	}
	if limit < 0 {
		return status.Error(codes.OutOfRange, "negative read limit")
	}

	file, size, err := s.binaryService.Open(stream.Context(), id)
	if err != nil {
		return s.contentError(err)
	}
	defer file.Close()

	if offset > size {
		return status.Error(codes.OutOfRange, "offset exceeds file size")
	}
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return s.contentError(err)
	}

	remaining := size - offset
	if limit > 0 && limit < remaining {
		remaining = limit
	}

	buf := make([]byte, 64*1024)
	for remaining > 0 {
		want := min(int64(len(buf)), remaining)
		n, err := io.ReadFull(file, buf[:int(want)])
		if err != nil {
			return s.contentError(err)
		}

		if err := stream.Send(&bytestream.ReadResponse{
			Data: bytes.Clone(buf[:n]),
		}); err != nil {
			return err
		}

		remaining -= int64(n)
	}

	return nil
}

func (s *ByteStreamServer) Write(stream bytestream.ByteStream_WriteServer) error {
	ctx := stream.Context()

	req, err := stream.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return status.Error(codes.InvalidArgument, "empty stream")
		}
		return err
	}

	resourceName := req.ResourceName

	id, err := s.parseID(resourceName)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	var committed int64

	for {
		if next := req.GetResourceName(); next != "" && next != resourceName {
			return status.Error(codes.InvalidArgument, "resource name changed")
		}
		if req.GetWriteOffset() < 0 {
			return status.Error(codes.InvalidArgument, "negative ")
		}

		committed, err = s.binaryService.WriteChunk(
			ctx, id, req.GetWriteOffset(), req.GetData(), req.GetFinishWrite(),
		)
		if err != nil {
			return s.contentError(err)
		}

		// Закрыть стрим без finish_write допустимо, поэтому
		// сначала проверим EOF
		finished := req.GetFinishWrite()

		req, err = stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return stream.SendAndClose(&bytestream.WriteResponse{
					CommittedSize: committed,
				})
			}
			return err
		}

		if finished {
			return status.Error(codes.InvalidArgument, "message received after finish_write")
		}
	}
}

func (s *ByteStreamServer) QueryWriteStatus(
	ctx context.Context, req *bytestream.QueryWriteStatusRequest,
) (*bytestream.QueryWriteStatusResponse, error) {
	id, err := s.parseID(req.GetResourceName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	size, complete, err := s.binaryService.WriteStatus(ctx, id)
	if err != nil {
		return nil, s.contentError(err)
	}

	return &bytestream.QueryWriteStatusResponse{
		CommittedSize: size,
		Complete:      complete,
	}, nil
}

func (s *ByteStreamServer) parseID(resourceName string) (int64, error) {
	raw, ok := strings.CutPrefix(resourceName, "binary/")
	if !ok {
		return 0, fmt.Errorf("invalid id")
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}

func (s *ByteStreamServer) contentError(err error) error {
	switch {
	case errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		return status.FromContextError(err).Err()

	case errors.Is(err, server.ErrDataNotFound):
		return status.Error(codes.NotFound, "binary not found")

	// TODO: описать ошибки

	// case errors.Is(err, server.ErrInvalidOffset):
	// 	return status.Error(codes.OutOfRange, "invalid write offset")

	// case errors.Is(err, server.ErrSizeMismatch):
	// 	return status.Error(codes.InvalidArgument, "size mismatch")

	case errors.Is(err, server.ErrUploadIncomplete):
		return status.Error(codes.FailedPrecondition, "upload incomplete")

	// case errors.Is(err, server.ErrUploadComplete):
	// 	return status.Error(codes.FailedPrecondition, "upload already complete")

	default:
		// Здесь также залогировать исходную ошибку.
		return status.Error(codes.Internal, "binary storage error")
	}
}
