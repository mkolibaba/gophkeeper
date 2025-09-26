package server

import (
	"context"
	"errors"
)

// TODO: какой порядок объявления?

var ErrDataAlreadyExists = errors.New("data already exists")
var ErrDataNotFound = errors.New("data not found")

type (
	Data interface {
		LoginData | NoteData | BinaryData | CardData
	}

	LoginData struct {
		User     string
		Name     string
		Login    string
		Password string
		Metadata map[string]string
	}

	NoteData struct {
		User     string
		Name     string
		Text     string
		Metadata map[string]string
	}

	BinaryData struct {
		User     string
		Name     string
		Data     []byte
		Metadata map[string]string
	}

	CardData struct {
		User       string
		Name       string
		Number     string
		ExpDate    string
		CVV        string
		Cardholder string
		Metadata   map[string]string
	}

	TypedDataService[T Data] interface {
		Save(ctx context.Context, data T) error
		GetAll(ctx context.Context, user string) ([]T, error)
		// TODO: update method
		Remove(ctx context.Context, name string, user string) error
	}

	LoginService interface {
		TypedDataService[LoginData]
	}

	NoteService interface {
		TypedDataService[NoteData]
	}

	BinaryService interface {
		TypedDataService[BinaryData]
	}

	CardService interface {
		TypedDataService[CardData]
	}
)
