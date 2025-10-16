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

type LoginServiceServer struct {
	gophkeeperv1.UnimplementedLoginServiceServer
	loginService server.LoginService
	validate     *validator.Validate
	logger       *log.Logger
}

func NewLoginServiceServer(
	loginService server.LoginService,
	validate *validator.Validate,
	logger *log.Logger,
) *LoginServiceServer {
	return &LoginServiceServer{
		loginService: loginService,
		validate:     validate,
		logger:       logger,
	}
}

func (s *LoginServiceServer) Save(ctx context.Context, in *gophkeeperv1.Login) (*empty.Empty, error) {
	data := server.LoginData{
		Name:     in.GetName(),
		Login:    in.GetLogin(),
		Password: in.GetPassword(),
		Website:  in.GetWebsite(),
		Notes:    in.GetNotes(),
	}

	if err := s.validate.StructCtx(ctx, &data); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.loginService.Create(ctx, data); err != nil {
		if errors.Is(err, server.ErrDataAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "data with this name already exists")
		}
		s.logger.Error("failed to save data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *LoginServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*gophkeeperv1.GetAllLoginsResponse, error) {
	logins, err := s.loginService.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve login data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var result []*gophkeeperv1.Login
	for _, login := range logins {
		var out gophkeeperv1.Login
		out.SetName(login.Name)
		out.SetLogin(login.Login)
		out.SetPassword(login.Password)
		out.SetWebsite(login.Website)
		out.SetNotes(login.Notes)
		result = append(result, &out)
	}

	var out gophkeeperv1.GetAllLoginsResponse
	out.SetResult(result)

	return &out, nil
}

func (s *LoginServiceServer) Remove(ctx context.Context, in *gophkeeperv1.RemoveDataRequest) (*empty.Empty, error) {
	if !in.HasId() {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.loginService.Remove(ctx, in.GetId()); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}
