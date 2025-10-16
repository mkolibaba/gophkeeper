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
	ErrPermissionDenied  = errors.New("permission denied")
)

type LoginData struct {
	ID       int64
	Name     string `validate:"required"`
	Login    string `validate:"required"`
	Password string
	Website  string
	Notes    string
	User     string
}

func (l LoginData) GetUser() string {
	return l.User
}

type LoginDataUpdate struct {
	Name     *string
	Login    *string
	Password *string
	Website  *string
	Notes    *string
}

// LoginService - сервис для работы с авторизационными данными типа логин/пароль.
type LoginService interface {
	// Create сохраняет данные для текущего пользователя.
	Create(ctx context.Context, data LoginData) error

	// GetAll возвращает все данные текущего пользователя.
	GetAll(ctx context.Context) ([]LoginData, error)

	// Update обновляет данные с переданным id. Только владелец
	// данных может редактировать их.
	Update(ctx context.Context, id int64, data LoginDataUpdate) error

	// Remove удаляет данные. Только владелец данных может
	// удалять их.
	Remove(ctx context.Context, id int64) error
}

type NoteData struct {
	ID   int64
	Name string `validate:"required"`
	Text string
	User string
}

func (n NoteData) GetUser() string {
	return n.User
}

type NoteDataUpdate struct {
	Text *string
}

// TODO(critical): для каждого сервиса добавить документацию, что манипулировать данными
//  может только владелец

// NoteService - сервис для работы с текстовыми данными.
type NoteService interface {
	// Create сохраняет текстовые данные для текущего пользователя.
	Create(ctx context.Context, data NoteData) error

	// GetAll возвращает все текстовые данные текущего пользователя.
	GetAll(ctx context.Context) ([]NoteData, error)

	// Update обновляет бинарные данные с переданным id. Только владелец
	// данных может редактировать их.
	Update(ctx context.Context, id int64, data NoteDataUpdate) error

	// Remove удаляет данные. Только владелец данных может
	// удалять их.
	Remove(ctx context.Context, id int64) error
}

type BinaryData struct {
	ID       int64
	Name     string `validate:"required"`
	Filename string `validate:"required"`
	Size     int64
	Notes    string
	User     string
}

func (b BinaryData) GetUser() string {
	return b.User
}

type ReadableBinaryData struct {
	BinaryData
	DataReader io.ReadCloser
}

type BinaryDataUpdate struct {
	Notes *string
}

// BinaryService - сервис для работы с бинарными данными.
type BinaryService interface {
	// Create сохраняет бинарные данные для текущего пользователя.
	Create(ctx context.Context, data ReadableBinaryData) error

	// Get возвращает бинарные данные с переданным id.
	Get(ctx context.Context, id int64) (*ReadableBinaryData, error)

	// GetAll возвращает все бинарные данные текущего пользователя.
	GetAll(ctx context.Context) ([]BinaryData, error)

	// Update обновляет бинарные данные с переданным id. Только владелец
	// данных может редактировать их.
	Update(ctx context.Context, id int64, data BinaryDataUpdate) error

	// Remove удаляет данные. Только владелец данных может
	// удалять их.
	Remove(ctx context.Context, id int64) error
}

type CardData struct {
	ID         int64
	Name       string `validate:"required"`
	Number     string `validate:"required,credit_card"`
	ExpDate    string `validate:"required,exp_date"`
	CVV        string `validate:"required,len=3"`
	Cardholder string `validate:"required"`
	Notes      string
	User       string
}

func (c CardData) GetUser() string {
	return c.User
}

type CardDataUpdate struct {
	Name       *string
	Number     *string
	ExpDate    *string
	CVV        *string
	Cardholder *string
	Notes      *string
}

// CardService - сервис для работы с данными карт.
type CardService interface {
	// Create сохраняет данные карты для текущего пользователя.
	Create(ctx context.Context, data CardData) error

	// GetAll возвращает все данные карт текущего пользователя.
	GetAll(ctx context.Context) ([]CardData, error)

	// Update обновляет данные карты с переданным id. Только владелец
	// данных может редактировать их.
	Update(ctx context.Context, id int64, data CardDataUpdate) error

	// Remove удаляет данные. Только владелец данных может
	// удалять их.
	Remove(ctx context.Context, id int64) error
}

type Data interface {
	GetUser() string
}

func VerifyCanEditData(ctx context.Context, data Data) error {
	if !IsCurrentUser(ctx, data.GetUser()) {
		return ErrPermissionDenied
	}
	return nil
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
