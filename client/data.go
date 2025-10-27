package client

import (
	"context"
	"github.com/go-playground/validator/v10"
	"regexp"
)

type Data interface {
	GetID() int64
	GetName() string
}

type LoginData struct {
	ID       int64
	Name     string `validate:"required"`
	Login    string `validate:"required"`
	Password string
	Website  string
	Notes    string
}

func (d LoginData) GetID() int64 {
	return d.ID
}

func (d LoginData) GetName() string {
	return d.Name
}

type LoginDataUpdate struct {
	ID       int64
	Name     *string
	Login    *string
	Password *string
	Website  *string
	Notes    *string
}

type LoginService interface {
	Save(ctx context.Context, data LoginData) error
	GetAll(ctx context.Context) ([]LoginData, error)
	Update(ctx context.Context, data LoginDataUpdate) error
	Remove(ctx context.Context, id int64) error
}

type NoteData struct {
	ID   int64
	Name string `validate:"required"`
	Text string
}

func (d NoteData) GetID() int64 {
	return d.ID
}

func (d NoteData) GetName() string {
	return d.Name
}

type NoteDataUpdate struct {
	ID   int64
	Name *string
	Text *string
}

type NoteService interface {
	Save(ctx context.Context, data NoteData) error
	GetAll(ctx context.Context) ([]NoteData, error)
	Update(ctx context.Context, data NoteDataUpdate) error
	Remove(ctx context.Context, id int64) error
}

type BinaryData struct {
	ID       int64
	Name     string `validate:"required"`
	Filename string `validate:"required"`
	Size     int64
	Notes    string
}

func (d BinaryData) GetID() int64 {
	return d.ID
}

func (d BinaryData) GetName() string {
	return d.Name
}

type BinaryDataUpdate struct {
	ID    int64
	Name  *string
	Notes *string
}

type BinaryService interface {
	Save(ctx context.Context, data BinaryData) error
	GetAll(ctx context.Context) ([]BinaryData, error)
	Update(ctx context.Context, data BinaryDataUpdate) error
	Remove(ctx context.Context, id int64) error
	Download(ctx context.Context, id int64) error
}

type CardData struct {
	ID         int64
	Name       string `validate:"required"`
	Number     string `validate:"required,credit_card"`
	ExpDate    string `validate:"required,exp_date"`
	CVV        string `validate:"required,len=3"`
	Cardholder string `validate:"required"`
	Notes      string
}

func (d CardData) GetID() int64 {
	return d.ID
}

func (d CardData) GetName() string {
	return d.Name
}

type CardDataUpdate struct {
	ID         int64
	Name       *string
	Number     *string
	ExpDate    *string
	CVV        *string
	Cardholder *string
	Notes      *string
}

type CardService interface {
	Save(ctx context.Context, data CardData) error
	GetAll(ctx context.Context) ([]CardData, error)
	Update(ctx context.Context, data CardDataUpdate) error
	Remove(ctx context.Context, id int64) error
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
