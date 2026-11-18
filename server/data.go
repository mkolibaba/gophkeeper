package server

import (
	"context"
	"errors"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	ErrDataNotFound      = errors.New("data not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidOffset     = errors.New("invalid offset")
	ErrUploadIncomplete  = errors.New("binary file upload has not been finished")
)

type LoginData struct {
	ID        int64
	Name      string `validate:"required"`
	Login     string `validate:"required"`
	Password  string
	Website   string
	Notes     string
	UpdatedAt time.Time
}

type LoginDataUpdate struct {
	Name     *string
	Login    *string
	Password *string
	Website  *string
	Notes    *string
}

// LoginService - сервис для работы с авторизационными данными типа логин/пароль.
// Работать с данными может только их владелец.
type LoginService interface {
	// Create сохраняет данные для текущего пользователя.
	Create(ctx context.Context, data LoginData) (int64, time.Time, error)

	// GetAll возвращает все данные текущего пользователя.
	GetAll(ctx context.Context) ([]LoginData, error)

	// Update обновляет данные с переданным id. Только владелец
	// данных может редактировать их.
	Update(ctx context.Context, id int64, data LoginDataUpdate) (time.Time, error)

	// Remove удаляет данные. Только владелец данных может
	// удалять их.
	Remove(ctx context.Context, id int64) error
}

type NoteData struct {
	ID        int64
	Name      string `validate:"required"`
	Text      string
	UpdatedAt time.Time
}

type NoteDataUpdate struct {
	Name *string
	Text *string
}

// NoteService - сервис для работы с текстовыми данными.
// Работать с данными может только их владелец.
type NoteService interface {
	// Create сохраняет текстовые данные для текущего пользователя.
	Create(ctx context.Context, data NoteData) (int64, time.Time, error)

	// GetAll возвращает все текстовые данные текущего пользователя.
	GetAll(ctx context.Context) ([]NoteData, error)

	// Update обновляет бинарные данные с переданным id. Только владелец
	// данных может редактировать их.
	Update(ctx context.Context, id int64, data NoteDataUpdate) (time.Time, error)

	// Remove удаляет данные. Только владелец данных может
	// удалять их.
	Remove(ctx context.Context, id int64) error
}

type BinaryData struct {
	ID        int64
	Name      string `validate:"required"`
	Filename  string `validate:"required"`
	Size      int64
	Notes     string
	UpdatedAt time.Time
}

type BinaryDataUpdate struct {
	Name  *string
	Notes *string
}

// BinaryService - сервис для работы с бинарными данными.
// Работать с данными может только их владелец.
type BinaryService interface {
	// Create сохраняет бинарные данные для текущего пользователя.
	Create(ctx context.Context, data BinaryData) (int64, time.Time, error)

	// Upload записывает содержимое ранее созданной записи.
	Upload(ctx context.Context, id int64, content io.Reader) error

	// GetAll возвращает все бинарные данные текущего пользователя.
	GetAll(ctx context.Context) ([]BinaryData, error)

	// Open открывает файл для чтения. 
	Open(ctx context.Context, id int64) (io.ReadSeekCloser, int64, error)

	// Update обновляет бинарные данные с переданным id. Только владелец
	// данных может редактировать их.
	Update(ctx context.Context, id int64, data BinaryDataUpdate) (time.Time, error)

	// Remove удаляет данные. Только владелец данных может удалять их.
	Remove(ctx context.Context, id int64) error

	// WriteChunk сохраняет чанк в pending-файл.
	WriteChunk(ctx context.Context, id int64, offset int64, data []byte, finish bool) (committedSize int64, err error)

	// WriteStatus возвращает сохранённый размер и состояние.
	// Если загрузка ещё не начиналась, возвращает ErrDataNotFound.
	WriteStatus(ctx context.Context, id int64) (committedSize int64, complete bool, err error)
}

type CardData struct {
	ID         int64
	Name       string `validate:"required"`
	Number     string `validate:"required,credit_card"`
	ExpDate    string `validate:"required,exp_date"`
	CVV        string `validate:"required,len=3"`
	Cardholder string `validate:"required"`
	Notes      string
	UpdatedAt  time.Time
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
// Работать с данными может только их владелец.
type CardService interface {
	// Create сохраняет данные карты для текущего пользователя.
	Create(ctx context.Context, data CardData) (int64, time.Time, error)

	// GetAll возвращает все данные карт текущего пользователя.
	GetAll(ctx context.Context) ([]CardData, error)

	// Update обновляет данные карты с переданным id. Только владелец
	// данных может редактировать их.
	Update(ctx context.Context, id int64, data CardDataUpdate) (time.Time, error)

	// Remove удаляет данные. Только владелец данных может
	// удалять их.
	Remove(ctx context.Context, id int64) error
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
