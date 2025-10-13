package grpc

import (
	"context"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"github.com/mkolibaba/gophkeeper/internal/server/grpc/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NoteServiceServer struct {
	pb.UnimplementedNoteServiceServer
	noteService   server.NoteService
	dataValidator *validator.Validate
	logger        *log.Logger
}

func NewNoteServiceServer(
	noteService server.NoteService,
	dataValidator *validator.Validate,
	logger *log.Logger,
) *NoteServiceServer {
	return &NoteServiceServer{
		noteService:   noteService,
		dataValidator: dataValidator,
		logger:        logger,
	}
}

func (s *NoteServiceServer) Save(ctx context.Context, in *pb.Note) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	data := server.NoteData{
		Name: in.GetName(),
		Text: in.GetText(),
	}

	if err := s.dataValidator.StructCtx(ctx, &data); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.noteService.Save(ctx, data, user); err != nil {
		if errors.Is(err, server.ErrDataAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "data with this name already exists")
		}
		s.logger.Error("failed to save data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *NoteServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*pb.GetAllNotesResponse, error) {
	user := utils.UserFromContext(ctx)

	notes, err := s.noteService.GetAll(ctx, user)

	if err != nil {
		s.logger.Error("failed to retrieve note data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var result []*pb.Note
	for _, note := range notes {
		var out pb.Note
		out.SetName(note.Name)
		out.SetText(note.Text)
		result = append(result, &out)
	}

	var out pb.GetAllNotesResponse
	out.SetResult(result)

	return &out, nil
}

func (s *NoteServiceServer) Remove(ctx context.Context, in *pb.RemoveDataRequest) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if err := s.noteService.Remove(ctx, in.GetName(), user); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}
