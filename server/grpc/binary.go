package grpc

import (
	"context"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"github.com/mkolibaba/gophkeeper/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
)

type BinaryServiceServer struct {
	gophkeeperv1.UnimplementedBinaryServiceServer
	binaryService server.BinaryService
	validate      *validator.Validate
	logger        *log.Logger
}

func NewBinaryServiceServer(
	binaryService server.BinaryService,
	validate *validator.Validate,
	logger *log.Logger,
) *BinaryServiceServer {
	return &BinaryServiceServer{
		binaryService: binaryService,
		validate:      validate, // TODO: использовать
		logger:        logger,
	}
}

func (s *BinaryServiceServer) Upload(stream grpc.ClientStreamingServer[gophkeeperv1.SaveBinaryRequest, empty.Empty]) error {
	// TODO: хорошей практикой было бы сначала проверить, что такой сущности
	//  нет, а потом сохранять чанки

	file, err := os.CreateTemp("", "*.tmp")
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer os.Remove(file.Name())

	var (
		filename    string
		size        int64
		name        string
		notes       string
		initialized bool
	)

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
			filename = in.GetFilename()
			size = in.GetSize()
			notes = in.GetNotes()
			initialized = true
		}

		_, err = file.Write(in.GetChunk().GetData())
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

	err = s.binaryService.Create(stream.Context(), server.ReadableBinaryData{
		BinaryData: server.BinaryData{
			Name:     name,
			Filename: filename,
			Size:     size,
			Notes:    notes,
		},
		DataReader: file,
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return stream.SendAndClose(&empty.Empty{})
}

func (s *BinaryServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*gophkeeperv1.GetAllBinariesResponse, error) {
	binaries, err := s.binaryService.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve binary data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var result []*gophkeeperv1.Binary
	for _, binary := range binaries {
		var out gophkeeperv1.Binary
		out.SetName(binary.Name)
		out.SetFilename(binary.Filename)
		out.SetSize(binary.Size)
		out.SetNotes(binary.Notes)
		result = append(result, &out)
	}

	var out gophkeeperv1.GetAllBinariesResponse
	out.SetResult(result)

	return &out, nil
}

func (s *BinaryServiceServer) Remove(ctx context.Context, in *gophkeeperv1.RemoveDataRequest) (*empty.Empty, error) {
	if !in.HasId() {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.binaryService.Remove(ctx, in.GetId()); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *BinaryServiceServer) Download(in *gophkeeperv1.DownloadBinaryRequest, stream grpc.ServerStreamingServer[gophkeeperv1.DownloadBinaryResponse]) error {
	binary, err := s.binaryService.Get(stream.Context(), in.GetId())
	if err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to retrieve data", "err", err)
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

		var chunk gophkeeperv1.FileChunk
		chunk.SetData(buf[:n])

		var out gophkeeperv1.DownloadBinaryResponse
		out.SetChunk(&chunk)
		out.SetName(binary.Name)
		out.SetFilename(binary.Filename)
		out.SetSize(binary.Size)

		if err := stream.Send(&out); err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		idx++
	}

	return nil
}
