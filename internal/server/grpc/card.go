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

type CardServiceServer struct {
	pb.UnimplementedCardServiceServer
	cardService   server.CardService
	dataValidator *validator.Validate
	logger        *zap.Logger
}

func NewCardServiceServer(
	cardService server.CardService,
	dataValidator *validator.Validate,
	logger *zap.Logger,
) *CardServiceServer {
	return &CardServiceServer{
		cardService:   cardService,
		dataValidator: dataValidator,
		logger:        logger,
	}
}

func (s *CardServiceServer) Save(ctx context.Context, in *pb.Card) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	data := server.CardData{
		User:       user,
		Name:       in.GetName(),
		Number:     in.GetNumber(),
		ExpDate:    in.GetExpDate(),
		CVV:        in.GetCvv(),
		Cardholder: in.GetCardholder(),
		Metadata:   in.GetMetadata(),
	}

	if err := s.dataValidator.StructCtx(ctx, &data); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.cardService.Save(ctx, data); err != nil {
		if errors.Is(err, server.ErrDataAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "data with this name already exists")
		}
		s.logger.Error("failed to save data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *CardServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*pb.GetAllCardsResponse, error) {
	user := utils.UserFromContext(ctx)

	cards, err := s.cardService.GetAll(ctx, user)
	if err != nil {
		s.logger.Error("failed to retrieve card data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var result []*pb.Card
	for _, card := range cards {
		var out pb.Card
		out.SetName(card.Name)
		out.SetNumber(card.Number)
		out.SetExpDate(card.ExpDate)
		out.SetCvv(card.CVV)
		out.SetCardholder(card.Cardholder)
		out.SetMetadata(card.Metadata)
		result = append(result, &out)
	}

	var out pb.GetAllCardsResponse
	out.SetResult(result)

	return &out, nil
}

func (s *CardServiceServer) Remove(ctx context.Context, in *pb.RemoveDataRequest) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if err := s.cardService.Remove(ctx, in.GetName(), user); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}
