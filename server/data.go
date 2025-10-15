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

// TODO(critical): поменять идентификатор name на id
type LoginData struct {
	Name     string `validate:"required"`
	Login    string `validate:"required"`
	Password string
	Website  string
	Notes    string
}

type LoginDataUpdate struct {
	Name     *string
	Login    *string
	Password *string
	Website  *string
	Notes    *string
}

type LoginService interface {
	Save(ctx context.Context, data LoginData, user string) error
	GetAll(ctx context.Context, user string) ([]LoginData, error)
	Update(ctx context.Context, data LoginDataUpdate, user string) error
	Remove(ctx context.Context, name string, user string) error
}

type NoteData struct {
	Name string `validate:"required"`
	Text string
}

type NoteDataUpdate struct {
	Text *string
}

type NoteService interface {
	Save(ctx context.Context, data NoteData, user string) error
	GetAll(ctx context.Context, user string) ([]NoteData, error)
	Update(ctx context.Context, data NoteDataUpdate, user string) error
	Remove(ctx context.Context, name string, user string) error
}

type BinaryData struct {
	Name     string `validate:"required"`
	Filename string `validate:"required"`
	Size     int64
	Notes    string
}

type ReadableBinaryData struct {
	BinaryData
	DataReader io.ReadCloser
}

type BinaryDataUpdate struct {
	Notes *string
}

type BinaryService interface {
	Save(ctx context.Context, data ReadableBinaryData, user string) error
	Get(ctx context.Context, name string, user string) (*ReadableBinaryData, error)
	GetAll(ctx context.Context, user string) ([]BinaryData, error)
	Update(ctx context.Context, data BinaryDataUpdate, user string) error
	Remove(ctx context.Context, name string, user string) error
}

type CardData struct {
	Name       string `validate:"required"`
	Number     string `validate:"required,credit_card"`
	ExpDate    string `validate:"required,exp_date"`
	CVV        string `validate:"required,len=3"`
	Cardholder string `validate:"required"`
	Notes      string
}

type CardService interface {
	Save(ctx context.Context, data CardData, user string) error
	GetAll(ctx context.Context, user string) ([]CardData, error)
	Update(ctx context.Context, data CardData, user string) error
	Remove(ctx context.Context, name string, user string) error
}

func RegisterDataValidationRules(validate *validator.Validate) error {
	expDateRegexp, err := regexp.Compile(`^\d{2}/\d{2}$`)
	if err != nil {
		return err
	}

	return validate.RegisterValidation("exp_date", func(fl validator.FieldLevel) bool {
		return expDateRegexp.MatchString(fl.Field().String())
	})
}
