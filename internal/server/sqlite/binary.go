package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/server"
	sqlc "github.com/mkolibaba/gophkeeper/internal/server/sqlite/sqlc/gen"
)

type BinaryService struct {
	qs *sqlc.Queries
}

func NewBinaryService(queries *sqlc.Queries) *BinaryService {
	return &BinaryService{
		qs: queries,
	}
}

func (b *BinaryService) Save(ctx context.Context, data server.BinaryData) error {
	metadata, err := json.Marshal(data.Metadata)
	if err != nil {
		return fmt.Errorf("save: invalid metadata: %w", err)
	}

	err = b.qs.SaveBinary(ctx, sqlc.SaveBinaryParams{
		Name:     data.Name,
		Data:     data.Data,
		Metadata: metadata,
		User:     data.User,
	})

	return tryUnwrapSaveError(err)
}

func (b *BinaryService) GetAll(ctx context.Context, user string) ([]server.BinaryData, error) {
	binaries, err := b.qs.GetAllBinaries(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.BinaryData
	for _, binary := range binaries {
		metadata, err := unmarshalMetadata(binary.Metadata)
		if err != nil {
			return nil, fmt.Errorf("get all: %w", err)
		}

		result = append(result, server.BinaryData{
			User:     binary.User,
			Name:     binary.Name,
			Data:     binary.Data,
			Metadata: metadata,
		})
	}

	return result, nil
}

func (b *BinaryService) Remove(ctx context.Context, name string, user string) error {
	n, err := b.qs.RemoveBinary(ctx, name)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if n == 0 {
		return server.ErrDataNotFound
	}
	return nil
}
