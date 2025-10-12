package client

import (
	"context"
	"github.com/go-playground/validator/v10"
	"regexp"
)

type (
	Data interface {
		GetName() string
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
		Name     string `validate:"required"`
		Filename string `validate:"required"`
		Size     int64
		Notes    string
	}

	CardData struct {
		Name       string `validate:"required"`
		Number     string `validate:"required,credit_card"`
		ExpDate    string `validate:"required,exp_date"`
		CVV        string `validate:"required,len=3"`
		Cardholder string `validate:"required"`
		Notes      string
	}

	BaseDataService[T Data] interface {
		Save(ctx context.Context, data T) error
		GetAll(ctx context.Context) ([]T, error)
		Remove(ctx context.Context, name string) error
	}

	LoginService interface {
		BaseDataService[LoginData]
	}

	NoteService interface {
		BaseDataService[NoteData]
	}

	BinaryService interface {
		BaseDataService[BinaryData]
		Download(ctx context.Context, name string) error
	}

	CardService interface {
		BaseDataService[CardData]
	}
)

func (d LoginData) GetName() string {
	return d.Name
}

func (d NoteData) GetName() string {
	return d.Name
}

func (d BinaryData) GetName() string {
	return d.Name
}

func (d CardData) GetName() string {
	return d.Name
}

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
