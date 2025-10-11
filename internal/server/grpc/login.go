package grpc

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"github.com/mkolibaba/gophkeeper/internal/server/grpc/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginServiceServer struct {
	pb.UnimplementedLoginServiceServer
	loginService  server.LoginService
	dataValidator *validator.Validate
	logger        *zap.Logger
}

func NewLoginServiceServer(
	loginService server.LoginService,
	dataValidator *validator.Validate,
	logger *zap.Logger,
) *LoginServiceServer {
	return &LoginServiceServer{
		loginService:  loginService,
		dataValidator: dataValidator,
		logger:        logger,
	}
}

func (s *LoginServiceServer) Save(ctx context.Context, in *pb.Login) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	data := server.LoginData{
		Name:     in.GetName(),
		Login:    in.GetLogin(),
		Password: in.GetPassword(),
		Website:  in.GetWebsite(),
		Notes:    in.GetNotes(),
	}

	if err := s.dataValidator.StructCtx(ctx, &data); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.loginService.Save(ctx, data, user); err != nil {
		if errors.Is(err, server.ErrDataAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "data with this name already exists")
		}
		s.logger.Error("failed to save data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *LoginServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*pb.GetAllLoginsResponse, error) {
	user := utils.UserFromContext(ctx)

	logins, err := s.loginService.GetAll(ctx, user)
	if err != nil {
		s.logger.Error("failed to retrieve login data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var result []*pb.Login
	for _, login := range logins {
		var out pb.Login
		out.SetName(login.Name)
		out.SetLogin(login.Login)
		out.SetPassword(login.Password)
		out.SetWebsite(login.Website)
		out.SetNotes(login.Notes)
		result = append(result, &out)
	}

	var out pb.GetAllLoginsResponse
	out.SetResult(result)

	return &out, nil
}

func (s *LoginServiceServer) Remove(ctx context.Context, in *pb.RemoveDataRequest) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if err := s.loginService.Remove(ctx, in.GetName(), user); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}
