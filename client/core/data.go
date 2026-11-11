package core

import (
	"context"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Item interface {
	GetID() string
	GetName() string
}

type ItemHeader struct {
	ID   string
	Name string `validate:"required"`
}

func (h ItemHeader) GetID() string {
	return h.ID
}

func (h ItemHeader) GetName() string {
	return h.Name
}

type LoginData struct {
	ItemHeader
	Login    string `validate:"required"`
	Password string
	Website  string
	Notes    string
}

type LoginDataUpdate struct {
	ID       string
	Name     *string
	Login    *string
	Password *string
	Website  *string
	Notes    *string
}

type NoteData struct {
	ItemHeader
	Text string
}

type NoteDataUpdate struct {
	ID   string
	Name *string
	Text *string
}

type BinaryData struct {
	ItemHeader
	Filename string `validate:"required"`
	Size     int64
	Notes    string
}

type BinaryDataUpdate struct {
	ID    string
	Name  *string
	Notes *string
}

type CardData struct {
	ItemHeader
	Number     string `validate:"required,credit_card"`
	ExpDate    string `validate:"required,exp_date"`
	CVV        string `validate:"required,len=3"`
	Cardholder string `validate:"required"`
	Notes      string
}

type CardDataUpdate struct {
	ID         string
	Name       *string
	Number     *string
	ExpDate    *string
	CVV        *string
	Cardholder *string
	Notes      *string
}

type VaultService interface {
	Get(ctx context.Context) ([]Item, error)

	CreateBinary(ctx context.Context, data BinaryData) (BinaryData, error)
	CreateCard(ctx context.Context, data CardData) (CardData, error)
	CreateLogin(ctx context.Context, data LoginData) (LoginData, error)
	CreateNote(ctx context.Context, data NoteData) (NoteData, error)

	UpdateNote(ctx context.Context, data NoteDataUpdate) (time.Time, error)
	UpdateLogin(ctx context.Context, data LoginDataUpdate) (time.Time, error)
	UpdateCard(ctx context.Context, data CardDataUpdate) (time.Time, error)
	UpdateBinary(ctx context.Context, data BinaryDataUpdate) (time.Time, error)

	DownloadBinary(ctx context.Context, id string) error

	DeleteItem(ctx context.Context, id string) error
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
