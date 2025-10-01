package grpc

import (
	"errors"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type BinaryServiceServer struct {
	pb.UnimplementedBinaryServiceServer
	binaryService server.BinaryService
	logger        *zap.Logger
}

func NewBinaryServiceServer(
	binaryService server.BinaryService,
	logger *zap.Logger,
) *BinaryServiceServer {
	return &BinaryServiceServer{
		binaryService: binaryService,
		logger:        logger,
	}
}

func (s *BinaryServiceServer) Download(in *pb.DownloadBinaryRequest, stream grpc.ServerStreamingServer[pb.FileChunk]) error {
	binary, err := s.binaryService.Get(stream.Context(), in.GetName(), "demo") // TODO: user из ctx
	if err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to retrieve data", zap.Error(err))
		return status.Error(codes.Internal, "internal server error")
	}

	file := binary.DataReader
	buf := make([]byte, 64*1024) // 64 KB
	idx := 0

	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		var chunk pb.FileChunk
		chunk.SetChunkData(buf[:n])
		chunk.SetFilename(binary.FileName)
		chunk.SetTotalSize(binary.Size)
		chunk.SetChunkIndex(int32(idx))

		if err := stream.Send(&chunk); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		idx++
	}

	return nil
}
