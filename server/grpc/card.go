package grpc

import (
	"context"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/grpc/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CardServiceServer struct {
	gophkeeperv1.UnimplementedCardServiceServer
	cardService   server.CardService
	dataValidator *validator.Validate
	logger        *log.Logger
}

func NewCardServiceServer(
	cardService server.CardService,
	dataValidator *validator.Validate,
	logger *log.Logger,
) *CardServiceServer {
	return &CardServiceServer{
		cardService:   cardService,
		dataValidator: dataValidator,
		logger:        logger,
	}
}

func (s *CardServiceServer) Save(ctx context.Context, in *gophkeeperv1.Card) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	data := server.CardData{
		Name:       in.GetName(),
		Number:     in.GetNumber(),
		ExpDate:    in.GetExpDate(),
		CVV:        in.GetCvv(),
		Cardholder: in.GetCardholder(),
		Notes:      in.GetNotes(),
	}

	if err := s.dataValidator.StructCtx(ctx, &data); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.cardService.Save(ctx, data, user); err != nil {
		if errors.Is(err, server.ErrDataAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "data with this name already exists")
		}
		s.logger.Error("failed to save data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *CardServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*gophkeeperv1.GetAllCardsResponse, error) {
	user := utils.UserFromContext(ctx)

	cards, err := s.cardService.GetAll(ctx, user)
	if err != nil {
		s.logger.Error("failed to retrieve card data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var result []*gophkeeperv1.Card
	for _, card := range cards {
		var out gophkeeperv1.Card
		out.SetName(card.Name)
		out.SetNumber(card.Number)
		out.SetExpDate(card.ExpDate)
		out.SetCvv(card.CVV)
		out.SetCardholder(card.Cardholder)
		out.SetNotes(card.Notes)
		result = append(result, &out)
	}

	var out gophkeeperv1.GetAllCardsResponse
	out.SetResult(result)

	return &out, nil
}

func (s *CardServiceServer) Remove(ctx context.Context, in *gophkeeperv1.RemoveDataRequest) (*empty.Empty, error) {
	user := utils.UserFromContext(ctx)

	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if err := s.cardService.Remove(ctx, in.GetName(), user); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}
