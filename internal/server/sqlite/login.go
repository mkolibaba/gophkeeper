package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/server"
	sqlc "github.com/mkolibaba/gophkeeper/internal/server/sqlite/sqlc/gen"
)

type LoginService struct {
	qs *sqlc.Queries
}

func NewLoginService(queries *sqlc.Queries) *LoginService {
	return &LoginService{
		qs: queries,
	}
}

func (l *LoginService) Save(ctx context.Context, data server.LoginData) error {
	metadata, err := json.Marshal(data.Metadata)
	if err != nil {
		return fmt.Errorf("save: invalid metadata: %w", err)
	}

	err = l.qs.SaveLogin(ctx, sqlc.SaveLoginParams{
		Name:     data.Name,
		Login:    data.Login,
		Password: &data.Password,
		Metadata: metadata,
		User:     data.User,
	})

	return tryUnwrapSaveError(err)
}

func (l *LoginService) GetAll(ctx context.Context, user string) ([]server.LoginData, error) {
	logins, err := l.qs.GetAllLogins(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.LoginData
	for _, login := range logins {
		metadata, err := unmarshalMetadata(login.Metadata)
		if err != nil {
			return nil, fmt.Errorf("get all: %w", err)
		}

		result = append(result, server.LoginData{
			User:     login.User,
			Name:     login.Name,
			Login:    login.Login,
			Password: *login.Password,
			Metadata: metadata,
		})
	}

	return result, nil
}

func (l *LoginService) Remove(ctx context.Context, name string, user string) error {
	n, err := l.qs.RemoveLogin(ctx, name)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if n == 0 {
		return server.ErrDataNotFound
	}
	return nil
}
