package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
)

var (
	ErrDataAlreadyExists  = errors.New("data already exists")
	ErrDataNotFound       = errors.New("data not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid login or password")
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

	User struct {
		Login    string
		Password string // TODO: hash
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

	UserService interface {
		Get(ctx context.Context, login string) (User, error)
		Save(ctx context.Context, user User) error
	}
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

// TODO: тут ли ему место? нужен ли интерфейс?
type AuthService struct {
	userService UserService
	logger      *log.Logger
}

func NewAuthService(userService UserService, logger *log.Logger) *AuthService {
	return &AuthService{
		userService: userService,
		logger:      logger,
	}
}

// TODO: по факту сервис не авторизует, а просто проверяет креды, поэтому название некорректное
func (s *AuthService) Authorize(ctx context.Context, login string, password string) error {
	u, err := s.userService.Get(ctx, login)
	if errors.Is(err, ErrUserNotFound) {
		return ErrInvalidCredentials
	}
	if err != nil {
		s.logger.Error("user get error", "err", err)
		return fmt.Errorf("internal error")
	}
	if password != u.Password {
		return ErrInvalidCredentials
	}

	return nil
}
