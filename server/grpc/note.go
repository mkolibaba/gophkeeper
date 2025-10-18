package grpc

import (
	"context"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"github.com/mkolibaba/gophkeeper/server"
	grpcgen "github.com/mkolibaba/gophkeeper/server/grpc/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NoteServiceServer struct {
	gophkeeperv1.UnimplementedNoteServiceServer
	noteService server.NoteService
	validate    *validator.Validate
	logger      *log.Logger
}

func NewNoteServiceServer(
	noteService server.NoteService,
	validate *validator.Validate,
	logger *log.Logger,
) *NoteServiceServer {
	return &NoteServiceServer{
		noteService: noteService,
		validate:    validate,
		logger:      logger,
	}
}

func (s *NoteServiceServer) Save(ctx context.Context, in *gophkeeperv1.Note) (*empty.Empty, error) {
	data := server.NoteData{
		Name: in.GetName(),
		Text: in.GetText(),
	}

	if err := s.validate.StructCtx(ctx, &data); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.noteService.Create(ctx, data); err != nil {
		if errors.Is(err, server.ErrDataAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "data with this name already exists")
		}
		s.logger.Error("failed to save data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *NoteServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*gophkeeperv1.GetAllNotesResponse, error) {
	notes, err := s.noteService.GetAll(ctx)

	if err != nil {
		s.logger.Error("failed to retrieve note data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var result []*gophkeeperv1.Note
	for _, note := range notes {
		var out gophkeeperv1.Note
		out.SetName(note.Name)
		out.SetText(note.Text)
		result = append(result, &out)
	}

	var out gophkeeperv1.GetAllNotesResponse
	out.SetResult(result)

	return &out, nil
}

func (s *NoteServiceServer) Update(ctx context.Context, in *gophkeeperv1.Note) (*empty.Empty, error) {
	return updateData(ctx, in, func(i *gophkeeperv1.Note) server.NoteDataUpdate {
		return grpcgen.MapNoteDataUpdate(i)
	}, s.noteService.Update, s.logger)
}

func (s *NoteServiceServer) Remove(ctx context.Context, in *gophkeeperv1.RemoveDataRequest) (*empty.Empty, error) {
	return removeData(ctx, in, s.noteService.Remove, s.logger)
}
