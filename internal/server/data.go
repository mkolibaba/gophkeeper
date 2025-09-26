package server

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// TODO: какой порядок объявления?

var ErrDataAlreadyExists = errors.New("data already exists")
var ErrDataNotFound = errors.New("data not found")

type (
	Data interface {
		LoginData | NoteData | BinaryData | CardData
	}

	LoginData struct {
		User     string `validate:"required"`
		Name     string `validate:"required"`
		Login    string `validate:"required"`
		Password string
		Metadata map[string]string
	}

	NoteData struct {
		User     string `validate:"required"`
		Name     string `validate:"required"`
		Text     string
		Metadata map[string]string
	}

	BinaryData struct {
		User     string `validate:"required"`
		Name     string `validate:"required"`
		Data     []byte `validate:"required"`
		Metadata map[string]string
	}

	CardData struct {
		User       string `validate:"required"`
		Name       string `validate:"required"`
		Number     string `validate:"required,credit_card"`
		ExpDate    string `validate:"required,exp_date"`
		CVV        string `validate:"required,len=3"`
		Cardholder string `validate:"required"`
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

func NewDataValidator() (*validator.Validate, error) {
	expDateRegexp, err := regexp.Compile(`^\d{2}/\d{2}$`)
	if err != nil {
		return nil, err
	}

	v := validator.New()
	err = v.RegisterValidation("exp_date", func(fl validator.FieldLevel) bool {
		return expDateRegexp.MatchString(fl.Field().String())
	})

	return v, err
}
