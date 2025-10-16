package grpc

import (
	"context"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"github.com/mkolibaba/gophkeeper/server"
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
	if !in.HasId() {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	var data server.NoteDataUpdate
	if in.HasName() {
		name := in.GetName()
		data.Name = &name
	}
	if in.HasText() {
		text := in.GetText()
		data.Text = &text
	}

	if err := s.noteService.Update(ctx, in.GetId(), data); err != nil {
		if errors.Is(err, server.ErrPermissionDenied) {
			return nil, status.Error(codes.PermissionDenied, server.ErrPermissionDenied.Error())
		}
		s.logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (s *NoteServiceServer) Remove(ctx context.Context, in *gophkeeperv1.RemoveDataRequest) (*empty.Empty, error) {
	if !in.HasId() {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.noteService.Remove(ctx, in.GetId()); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}
