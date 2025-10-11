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
		Website  string
		Notes    string
	}

	NoteData struct {
		Data
		Name string `validate:"required"`
		Text string
	}

	BinaryData struct {
		Data
		Name     string `validate:"required"`
		Filename string `validate:"required"`
		Size     int64
		Notes    string
	}

	CardData struct {
		Data
		Name       string `validate:"required"`
		Number     string `validate:"required,credit_card"`
		ExpDate    string `validate:"required,exp_date"`
		CVV        string `validate:"required,len=3"`
		Cardholder string `validate:"required"`
		Notes      string
	}

	User struct {
		Login    string
		Password string
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

	// TODO: отрефакторить сервис
	AuthorizationService interface {
		Authorize(ctx context.Context, login string, password string) (string, error)
		Register(ctx context.Context, login string, password string) (string, error)
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

// TODO: подумать
type Session struct {
	user *User
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) SetCurrentUser(user User) {
	s.user = &user
}

func (s *Session) GetCurrentUser() *User {
	return s.user
}
