package grpc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"charm.land/log/v2"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type VaultServiceServer struct {
	gophkeeperv1.UnimplementedVaultServiceServer
	binaryService server.BinaryService
	cardService   server.CardService
	loginService  server.LoginService
	noteService   server.NoteService
	validate      *validator.Validate
	updaters      map[string]func(context.Context, int64, *gophkeeperv1.UpdateItemRequest) (time.Time, error)
	removers      map[string]func(context.Context, int64) error
	logger        *log.Logger
}

func NewVaultServiceServer(
	binaryService server.BinaryService,
	cardService server.CardService,
	loginService server.LoginService,
	noteService server.NoteService,
	validate *validator.Validate,
	logger *log.Logger,
) *VaultServiceServer {
	return &VaultServiceServer{
		binaryService: binaryService,
		cardService:   cardService,
		loginService:  loginService,
		noteService:   noteService,
		validate:      validate,
		updaters: map[string]func(context.Context, int64, *gophkeeperv1.UpdateItemRequest) (time.Time, error){
			"binary": func(ctx context.Context, id int64, item *gophkeeperv1.UpdateItemRequest) (time.Time, error) {
				var in server.BinaryDataUpdate
				if item.HasName() {
					in.Name = new(item.GetName())
				}
				if item.HasNotes() {
					in.Notes = new(item.GetNotes())
				}
				return binaryService.Update(ctx, id, in)
			},
			"card": func(ctx context.Context, id int64, item *gophkeeperv1.UpdateItemRequest) (time.Time, error) {
				var in server.CardDataUpdate
				if item.HasName() {
					in.Name = new(item.GetName())
				}
				if item.HasNotes() {
					in.Notes = new(item.GetNotes())
				}
				if item.HasCard() {
					card := item.GetCard()
					if card.HasCardholder() {
						in.Cardholder = new(card.GetCardholder())
					}
					if card.HasCvv() {
						in.CVV = new(card.GetCvv())
					}
					if card.HasExpDate() {
						in.ExpDate = new(card.GetExpDate())
					}
					if card.HasNumber() {
						in.Number = new(card.GetNumber())
					}
				}
				return cardService.Update(ctx, id, in)
			},
			"login": func(ctx context.Context, id int64, item *gophkeeperv1.UpdateItemRequest) (time.Time, error) {
				var in server.LoginDataUpdate
				if item.HasName() {
					in.Name = new(item.GetName())
				}
				if item.HasNotes() {
					in.Notes = new(item.GetNotes())
				}
				if item.HasLogin() {
					login := item.GetLogin()
					if login.HasLogin() {
						in.Login = new(login.GetLogin())
					}
					if login.HasPassword() {
						in.Password = new(login.GetPassword())
					}
					if login.HasWebsite() {
						in.Website = new(login.GetWebsite())
					}
				}
				return loginService.Update(ctx, id, in)
			},
			"note": func(ctx context.Context, id int64, item *gophkeeperv1.UpdateItemRequest) (time.Time, error) {
				var in server.NoteDataUpdate
				if item.HasName() {
					in.Name = new(item.GetName())
				}
				if item.HasNotes() {
					in.Text = new(item.GetNotes())
				}
				return noteService.Update(ctx, id, in)
			},
		},
		removers: map[string]func(context.Context, int64) error{
			"binary": binaryService.Remove,
			"card":   cardService.Remove,
			"login":  loginService.Remove,
			"note":   noteService.Remove,
		},
		logger: logger,
	}
}

func (s *VaultServiceServer) Get(ctx context.Context, _ *emptypb.Empty) (*gophkeeperv1.Vault, error) {
	var items []*gophkeeperv1.Item

	binaries, err := s.binaryService.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve binary data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	for _, b := range binaries {
		var item gophkeeperv1.Item
		item.SetId("binary/" + strconv.FormatInt(b.ID, 10))
		item.SetName(b.Name)
		item.SetNotes(b.Notes)
		item.SetUpdatedAt(timestamppb.New(b.UpdatedAt))
		var binary gophkeeperv1.Binary
		binary.SetFilename(b.Filename)
		binary.SetSize(b.Size)
		item.SetBinary(&binary)
		items = append(items, &item)
	}

	cards, err := s.cardService.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve card data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	for _, c := range cards {
		var item gophkeeperv1.Item
		item.SetId("card/" + strconv.FormatInt(c.ID, 10))
		item.SetName(c.Name)
		item.SetNotes(c.Notes)
		item.SetUpdatedAt(timestamppb.New(c.UpdatedAt))
		var card gophkeeperv1.Card
		card.SetNumber(c.Number)
		card.SetExpDate(c.ExpDate)
		card.SetCvv(c.CVV)
		card.SetCardholder(c.Cardholder)
		item.SetCard(&card)
		items = append(items, &item)
	}

	logins, err := s.loginService.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve login data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	for _, l := range logins {
		var item gophkeeperv1.Item
		item.SetId("login/" + strconv.FormatInt(l.ID, 10))
		item.SetName(l.Name)
		item.SetNotes(l.Notes)
		item.SetUpdatedAt(timestamppb.New(l.UpdatedAt))
		var login gophkeeperv1.Login
		login.SetLogin(l.Login)
		login.SetPassword(l.Password)
		login.SetWebsite(l.Website)
		item.SetLogin(&login)
		items = append(items, &item)
	}

	notes, err := s.noteService.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve note data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}
	for _, n := range notes {
		var item gophkeeperv1.Item
		item.SetId("note/" + strconv.FormatInt(n.ID, 10))
		item.SetName(n.Name)
		item.SetNotes(n.Text)
		item.SetUpdatedAt(timestamppb.New(n.UpdatedAt))
		item.SetNote(&gophkeeperv1.Note{})
		items = append(items, &item)
	}

	var out gophkeeperv1.Vault
	out.SetItems(items)
	return &out, nil
}

func (s *VaultServiceServer) CreateItem(ctx context.Context, in *gophkeeperv1.CreateItemRequest) (*gophkeeperv1.Item, error) {
	var out gophkeeperv1.Item
	out.SetName(in.GetName())
	out.SetNotes(in.GetNotes())

	switch {
	case in.HasBinary():
		data := server.BinaryData{
			Name:     in.GetName(),
			Filename: in.GetBinary().GetFilename(),
			Size:     in.GetBinary().GetSize(),
			Notes:    in.GetNotes(),
		}
		if err := s.validate.StructCtx(ctx, &data); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		id, updatedAt, err := s.binaryService.Create(ctx, data)
		if err != nil {
			s.logger.Error("failed to save data", "err", err)
			return nil, status.Error(codes.Internal, "internal server error")
		}

		out.SetId("binary/" + strconv.FormatInt(id, 10))
		out.SetUpdatedAt(timestamppb.New(updatedAt))
		out.SetBinary(in.GetBinary())
	case in.HasCard():
		data := server.CardData{
			Name:       in.GetName(),
			Number:     in.GetCard().GetNumber(),
			ExpDate:    in.GetCard().GetExpDate(),
			CVV:        in.GetCard().GetCvv(),
			Cardholder: in.GetCard().GetCardholder(),
			Notes:      in.GetNotes(),
		}
		if err := s.validate.StructCtx(ctx, &data); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		id, updatedAt, err := s.cardService.Create(ctx, data)
		if err != nil {
			s.logger.Error("failed to save data", "err", err)
			return nil, status.Error(codes.Internal, "internal server error")
		}

		out.SetId("card/" + strconv.FormatInt(id, 10))
		out.SetUpdatedAt(timestamppb.New(updatedAt))
		out.SetCard(in.GetCard())
	case in.HasLogin():
		data := server.LoginData{
			Name:     in.GetName(),
			Login:    in.GetLogin().GetLogin(),
			Password: in.GetLogin().GetPassword(),
			Website:  in.GetLogin().GetWebsite(),
			Notes:    in.GetNotes(),
		}

		if err := s.validate.StructCtx(ctx, &data); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		id, updatedAt, err := s.loginService.Create(ctx, data)
		if err != nil {
			s.logger.Error("failed to save data", "err", err)
			return nil, status.Error(codes.Internal, "internal server error")
		}

		out.SetId("login/" + strconv.FormatInt(id, 10))
		out.SetUpdatedAt(timestamppb.New(updatedAt))
		out.SetLogin(in.GetLogin())
	case in.HasNote():
		data := server.NoteData{
			Name: in.GetName(),
			Text: in.GetNotes(),
		}

		if err := s.validate.StructCtx(ctx, &data); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		id, updatedAt, err := s.noteService.Create(ctx, data)
		if err != nil {
			s.logger.Error("failed to save data", "err", err)
			return nil, status.Error(codes.Internal, "internal server error")
		}

		out.SetId("note/" + strconv.FormatInt(id, 10))
		out.SetUpdatedAt(timestamppb.New(updatedAt))
		out.SetNote(in.GetNote())
	}

	return &out, nil
}

func (s *VaultServiceServer) UpdateItem(ctx context.Context, in *gophkeeperv1.UpdateItemRequest) (*gophkeeperv1.Item, error) {
	id, typ, err := s.parseID(in.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	update, ok := s.updaters[typ]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	updatedAt, err := update(ctx, id, in)
	if err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to update data", "err", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var out gophkeeperv1.Item
	out.SetId(in.GetId())
	out.SetUpdatedAt(timestamppb.New(updatedAt))
	if in.HasName() {
		out.SetName(in.GetName())
	}
	if in.HasNotes() {
		out.SetNotes(in.GetNotes())
	}
	if in.HasBinary() {
		out.SetBinary(in.GetBinary())
	}
	if in.HasCard() {
		out.SetCard(in.GetCard())
	}
	if in.HasLogin() {
		out.SetLogin(in.GetLogin())
	}
	if in.HasNote() {
		out.SetNote(in.GetNote())
	}
	return &out, nil
}

func (s *VaultServiceServer) DeleteItem(ctx context.Context, in *gophkeeperv1.DeleteItemRequest) (*emptypb.Empty, error) {
	id, typ, err := s.parseID(in.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	remove, ok := s.removers[typ]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	if err := remove(ctx, id); err != nil {
		if errors.Is(err, server.ErrDataNotFound) {
			return nil, status.Error(codes.NotFound, "data not found")
		}
		s.logger.Error("failed to remove data", "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &empty.Empty{}, nil
}

func (s *VaultServiceServer) parseID(id string) (int64, string, error) {
	typ, idStr, found := strings.Cut(id, "/")
	if !found {
		return 0, "", fmt.Errorf("invalid id")
	}

	result, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid id")
	}

	return result, typ, nil
}
