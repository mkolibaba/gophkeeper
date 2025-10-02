package client

import (
	"context"
	"github.com/go-playground/validator/v10"
	"regexp"
)

type (
	Data interface {
	}

	LoginData struct {
		Data
		Name     string `validate:"required"`
		Login    string `validate:"required"`
		Password string
		Metadata map[string]string
	}

	NoteData struct {
		Data
		Name     string `validate:"required"`
		Text     string
		Metadata map[string]string
	}

	BinaryData struct {
		Data
		Name  string `validate:"required"`
		Bytes []byte `validate:"required"`
		//Size     int
		FileName string `validate:"required"`
		Metadata map[string]string
	}

	CardData struct {
		Data
		Name       string `validate:"required"`
		Number     string `validate:"required,credit_card"`
		ExpDate    string `validate:"required,exp_date"`
		CVV        string `validate:"required,len=3"`
		Cardholder string `validate:"required"`
		Metadata   map[string]string
	}

	TypedDataService[T Data] interface {
		Save(ctx context.Context, user string, data T) error
		GetAll(ctx context.Context, user string) ([]T, error)
		// TODO: update method
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
