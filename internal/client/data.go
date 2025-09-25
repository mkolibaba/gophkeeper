package client

import (
	"context"
	"errors"
)

const (
	DataTypeLogin DataType = iota
	DataTypeNote
	DataTypeBinary
	DataTypeCard
)

var (
	ErrEmptyName = errors.New("name is empty")
)

type (
	DataType uint8

	Data interface {
		Type() DataType
		Validate() error
	}

	LoginData struct {
		Name     string
		Login    string
		Password string
		Metadata map[string]string
	}

	NoteData struct {
		Name     string
		Text     string
		Metadata map[string]string
	}

	BinaryData struct {
		Name     string
		Data     []byte
		Metadata map[string]string
	}

	CardData struct {
		Name       string
		Number     string
		ExpDate    string
		CVV        string
		Cardholder string
		Metadata   map[string]string
	}
)

func (l LoginData) Type() DataType {
	return DataTypeLogin
}

func (l LoginData) Validate() error {
	if l.Name == "" {
		return ErrEmptyName
	}
	if l.Login == "" {
		return errors.New("login is empty")
	}
	return nil
}

func (n NoteData) Type() DataType {
	return DataTypeNote
}

func (n NoteData) Validate() error {
	if n.Name == "" {
		return ErrEmptyName
	}
	if n.Text == "" {
		return errors.New("text is empty")
	}
	return nil
}

func (b BinaryData) Type() DataType {
	return DataTypeBinary
}

func (b BinaryData) Validate() error {
	if b.Name == "" {
		return ErrEmptyName
	}
	if len(b.Data) == 0 {
		return errors.New("data is empty")
	}
	return nil
}

func (c CardData) Type() DataType {
	return DataTypeCard
}

func (c CardData) Validate() error {
	if c.Name == "" {
		return ErrEmptyName
	}
	if c.Number == "" {
		return errors.New("number is empty")
	}
	if c.ExpDate == "" {
		return errors.New("exp date is empty")
	}
	if c.CVV == "" {
		return errors.New("cvv is empty")
	}
	if c.Cardholder == "" {
		return errors.New("cardholder is empty")
	}
	return nil
}

type DataService interface {
	Save(ctx context.Context, user string, data Data) error
	GetAll(ctx context.Context, user string, dataType DataType) ([]Data, error)
}
