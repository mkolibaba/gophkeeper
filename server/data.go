package server

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
)

var (
	ErrDataAlreadyExists = errors.New("data already exists")
	ErrDataNotFound      = errors.New("data not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type (
	Data interface {
		LoginData | NoteData | BinaryData | CardData
	}

	LoginData struct {
		Name     string `validate:"required"`
		Login    string `validate:"required"`
		Password string
		Website  string
		Notes    string
	}

	NoteData struct {
		Name string `validate:"required"`
		Text string
	}

	BinaryData struct {
		Name       string `validate:"required"`
		Filename   string
		Size       int64
		DataReader io.ReadCloser // TODO: этот reader должен быть только у реквеста
		Notes      string
	}

	CardData struct {
		Name       string `validate:"required"`
		Number     string `validate:"required,credit_card"`
		ExpDate    string `validate:"required,exp_date"`
		CVV        string `validate:"required,len=3"`
		Cardholder string `validate:"required"`
		Notes      string
	}

	TypedDataService[T Data] interface {
		Save(ctx context.Context, data T, user string) error
		GetAll(ctx context.Context, user string) ([]T, error)
		// TODO: update method

		// TODO: тут нужен user? возможно, да, но только для валидации
		Remove(ctx context.Context, name string, user string) error
	}

	LoginService TypedDataService[LoginData]

	NoteService TypedDataService[NoteData]

	BinaryService interface {
		TypedDataService[BinaryData]
		Get(ctx context.Context, name string, user string) (BinaryData, error)
	}

	CardService TypedDataService[CardData]
)

// TODO: нужно ли разделить data validator и user validator? либо нормально их объединить?
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
