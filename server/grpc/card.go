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

type CardServiceServer struct {
	gophkeeperv1.UnimplementedCardServiceServer
	cardService server.CardService
	validate    *validator.Validate
	logger      *log.Logger
}

func NewCardServiceServer(
	cardService server.CardService,
	validate *validator.Validate,
	logger *log.Logger,
) *CardServiceServer {
	return &CardServiceServer{
		cardService: cardService,
		validate:    validate,
		logger:      logger,
	}
}

func (s *CardServiceServer) Save(ctx context.Context, in *gophkeeperv1.Card) (*empty.Empty, error) {
	data := server.CardData{
		Name:       in.GetName(),
		Number:     in.GetNumber(),
		ExpDate:    in.GetExpDate(),
		CVV:        in.GetCvv(),
		Cardholder: in.GetCardholder(),
		Notes:      in.GetNotes(),
	}

	if err := s.validate.StructCtx(ctx, &data); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.cardService.Create(ctx, data); err != nil {
		if errors.Is(err, server.ErrDataAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "data with this name already exists")
		}
		s.logger.Error("failed to save data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *CardServiceServer) GetAll(ctx context.Context, _ *empty.Empty) (*gophkeeperv1.GetAllCardsResponse, error) {
	cards, err := s.cardService.GetAll(ctx)
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

func (s *CardServiceServer) Update(ctx context.Context, in *gophkeeperv1.Card) (*empty.Empty, error) {
	return updateData(ctx, in, func(i *gophkeeperv1.Card) server.CardDataUpdate {
		return grpcgen.MapCardDataUpdate(i)
	}, s.cardService.Update, s.logger)
}

func (s *CardServiceServer) Remove(ctx context.Context, in *gophkeeperv1.RemoveDataRequest) (*empty.Empty, error) {
	return removeData(ctx, in, s.cardService.Remove, s.logger)
}
