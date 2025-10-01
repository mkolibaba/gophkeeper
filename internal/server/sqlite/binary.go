package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/server"
	sqlc "github.com/mkolibaba/gophkeeper/internal/server/sqlite/sqlc/gen"
	"io"
	"os"
	"path/filepath"
)

type BinaryService struct {
	qs             *sqlc.Queries
	binariesFolder string
}

func NewBinaryService(queries *sqlc.Queries, db *DB) *BinaryService {
	return &BinaryService{
		qs:             queries,
		binariesFolder: db.binariesFolder,
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

func (b *BinaryService) Get(ctx context.Context, name string, user string) (server.BinaryData, error) {
	binary, err := b.qs.GetBinary(ctx, name, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server.BinaryData{}, server.ErrDataNotFound
		}
		return server.BinaryData{}, fmt.Errorf("get: %w", err)
	}

	metadata, err := unmarshalMetadata(binary.Metadata)
	if err != nil {
		return server.BinaryData{}, fmt.Errorf("get: %w", err)
	}

	file, err := os.ReadFile(filepath.Join(b.binariesFolder, binary.Path)) // TODO: заменить на open
	if err != nil {
		return server.BinaryData{}, fmt.Errorf("get: %w", err)
	}

	return server.BinaryData{
		User:       binary.User,
		Name:       binary.Name,
		Data:       file,
		Metadata:   metadata,
		FileName:   binary.Path,
		DataReader: io.NopCloser(bytes.NewReader(file)),
		Size:       int64(len(file)),
	}, nil
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
			FileName: binary.Path,
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
