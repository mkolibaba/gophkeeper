package grpc

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dataServiceServer struct {
	pb.UnimplementedDataServiceServer
	loginService  server.LoginService
	noteService   server.NoteService
	binaryService server.BinaryService
	cardService   server.CardService
	dataValidator *validator.Validate
	logger        *zap.Logger
}

func (d *dataServiceServer) Save(ctx context.Context, in *pb.SaveDataRequest) (*empty.Empty, error) {
	var validate func() error // TODO: придумать решение получше
	var save func() error

	switch in.WhichData() {
	case pb.SaveDataRequest_Login_case:
		login := in.GetLogin()
		data := server.LoginData{
			User:     in.GetUser(),
			Name:     login.GetName(),
			Login:    login.GetLogin(),
			Password: login.GetPassword(),
			Metadata: login.GetMetadata(),
		}
		validate = func() error {
			return d.dataValidator.StructCtx(ctx, &data)
		}
		save = func() error {
			return d.loginService.Save(ctx, data)
		}
	case pb.SaveDataRequest_Note_case:
		note := in.GetNote()
		data := server.NoteData{
			User:     in.GetUser(),
			Name:     note.GetName(),
			Text:     note.GetText(),
			Metadata: note.GetMetadata(),
		}
		validate = func() error {
			return d.dataValidator.StructCtx(ctx, &data)
		}
		save = func() error {
			return d.noteService.Save(ctx, data)
		}
	case pb.SaveDataRequest_Binary_case:
		binary := in.GetBinary()
		data := server.BinaryData{
			User:     in.GetUser(),
			Name:     binary.GetName(),
			Data:     binary.GetData(),
			Metadata: binary.GetMetadata(),
		}
		validate = func() error {
			return d.dataValidator.StructCtx(ctx, &data)
		}
		save = func() error {
			return d.binaryService.Save(ctx, data)
		}
	case pb.SaveDataRequest_Card_case:
		card := in.GetCard()
		data := server.CardData{
			User:       in.GetUser(),
			Name:       card.GetName(),
			Number:     card.GetNumber(),
			ExpDate:    card.GetExpDate(),
			CVV:        card.GetCvv(),
			Cardholder: card.GetCardholder(),
			Metadata:   card.GetMetadata(),
		}
		validate = func() error {
			return d.dataValidator.StructCtx(ctx, &data)
		}
		save = func() error {
			return d.cardService.Save(ctx, data)
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid data type")
	}

	if err := validate(); err != nil {
		// TODO: посмотреть какая структура ошибки и сделать лучше
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := save(); err != nil {
		if errors.Is(err, server.ErrDataAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "data with this name already exists")
		}
		d.logger.Error("failed to save data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (d *dataServiceServer) GetAll(ctx context.Context, in *pb.GetAllDataRequest) (*pb.GetAllDataResponse, error) {
	if in.GetUser() == "" {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}

	var data []*pb.DataWrapper

	switch in.GetDataType() {
	case pb.DataType_LOGIN:
		logins, err := d.loginService.GetAll(ctx, in.GetUser())
		if err != nil {
			d.logger.Error("failed to retrieve data", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal server error")
		}

		for _, login := range logins {
			var out pb.Login
			out.SetName(login.Name)
			out.SetLogin(login.Login)
			out.SetPassword(login.Password)
			out.SetMetadata(login.Metadata)

			var wrapper pb.DataWrapper
			wrapper.SetLogin(&out)

			data = append(data, &wrapper)
		}
	case pb.DataType_NOTE:
		notes, err := d.noteService.GetAll(ctx, in.GetUser())
		if err != nil {
			d.logger.Error("failed to retrieve data", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal server error")
		}

		for _, note := range notes {
			var out pb.Note
			out.SetName(note.Name)
			out.SetText(note.Text)
			out.SetMetadata(note.Metadata)

			var wrapper pb.DataWrapper
			wrapper.SetNote(&out)

			data = append(data, &wrapper)
		}
	case pb.DataType_BINARY:
		binaries, err := d.binaryService.GetAll(ctx, in.GetUser())
		if err != nil {
			d.logger.Error("failed to retrieve data", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal server error")
		}

		for _, binary := range binaries {
			var out pb.Binary
			out.SetName(binary.Name)
			out.SetData(binary.Data)
			out.SetMetadata(binary.Metadata)

			var wrapper pb.DataWrapper
			wrapper.SetBinary(&out)

			data = append(data, &wrapper)
		}
	case pb.DataType_CARD:
		cards, err := d.cardService.GetAll(ctx, in.GetUser())
		if err != nil {
			d.logger.Error("failed to retrieve data", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal server error")
		}

		for _, card := range cards {
			var out pb.Card
			out.SetName(card.Name)
			out.SetNumber(card.Number)
			out.SetExpDate(card.ExpDate)
			out.SetCvv(card.CVV)
			out.SetCardholder(card.Cardholder)
			out.SetMetadata(card.Metadata)

			var wrapper pb.DataWrapper
			wrapper.SetCard(&out)

			data = append(data, &wrapper)
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid data type")
	}

	var out pb.GetAllDataResponse
	out.SetData(data)
	return &out, nil
}

func (d *dataServiceServer) Remove(ctx context.Context, in *pb.RemoveDataRequest) (*empty.Empty, error) {
	// TODO: validate name and user
	if in.GetUser() == "" {
		return nil, status.Error(codes.InvalidArgument, "user is required")
	}
	if in.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	var remove func() error

	switch in.GetDataType() {
	case pb.DataType_LOGIN:
		remove = func() error {
			return d.loginService.Remove(ctx, in.GetName(), in.GetUser())
		}
	case pb.DataType_NOTE:
		remove = func() error {
			return d.noteService.Remove(ctx, in.GetName(), in.GetUser())
		}
	case pb.DataType_BINARY:
		remove = func() error {
			return d.binaryService.Remove(ctx, in.GetName(), in.GetUser())
		}
	case pb.DataType_CARD:
		remove = func() error {
			return d.cardService.Remove(ctx, in.GetName(), in.GetUser())
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid data type")
	}

	if err := remove(); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		d.logger.Error("failed to remove data", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}
