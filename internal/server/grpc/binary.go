package grpc

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"github.com/mkolibaba/gophkeeper/internal/server/grpc/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
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

func (s *BinaryServiceServer) Upload(stream grpc.ClientStreamingServer[pb.SaveBinaryRequest, empty.Empty]) error {
	user := utils.UserFromContext(stream.Context())

	file, err := os.CreateTemp("", "*.tmp")
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer os.Remove(file.Name())

	var filename string
	var size int64
	var name string
	var metadata map[string]string
	var initialized bool

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if !initialized {
			name = in.GetName()
			filename = in.GetChunk().GetFilename()
			size = in.GetChunk().GetTotalSize()
			metadata = in.GetMetadata()
			initialized = true
		}

		_, err = file.Write(in.GetChunk().GetChunkData())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	if err = file.Sync(); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	err = s.binaryService.Save(stream.Context(), server.BinaryData{
		User:       user,
		Name:       name,
		FileName:   filename,
		DataReader: file,
		Size:       size,
		Metadata:   metadata,
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return stream.SendAndClose(&empty.Empty{})
}

func (s *BinaryServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*pb.GetAllBinariesResponse, error) {
	user := utils.UserFromContext(ctx)

	binaries, err := s.binaryService.GetAll(ctx, user)
	if err != nil {
		s.logger.Error("failed to retrieve binary data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var result []*pb.Binary
	for _, binary := range binaries {
		var out pb.Binary
		out.SetName(binary.Name)
		out.SetFileName(binary.FileName)
		out.SetMetadata(binary.Metadata)
		result = append(result, &out)
	}

	var out pb.GetAllBinariesResponse
	out.SetResult(result)

	return &out, nil
}

func (s *BinaryServiceServer) Remove(ctx context.Context, in *pb.RemoveDataRequest) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	// TODO: удалить файл тоже

	if err := s.binaryService.Remove(ctx, in.GetName(), user); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *BinaryServiceServer) Download(in *pb.DownloadBinaryRequest, stream grpc.ServerStreamingServer[pb.FileChunk]) error {
	user := utils.UserFromContext(stream.Context())

	binary, err := s.binaryService.Get(stream.Context(), in.GetName(), user)
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
