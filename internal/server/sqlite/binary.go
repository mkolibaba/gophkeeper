package sqlite

import (
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

	// TODO: тут лучше проверить, что такой записи и такого файла нет

	err = b.qs.SaveBinary(ctx, sqlc.SaveBinaryParams{
		Name:     data.Name,
		Filename: data.FileName,
		Metadata: metadata,
		User:     data.User,
	})

	if err != nil {
		return tryUnwrapSaveError(err)
	}

	dest, err := os.Create(filepath.Join(b.binariesFolder, data.Name))
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}
	defer dest.Close()

	buf := make([]byte, 1024*1024) // 1 MB
	_, err = io.CopyBuffer(dest, data.DataReader, buf)

	if err != nil {
		os.Remove(dest.Name())
		return fmt.Errorf("save: %w", err)
	}

	return nil
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

	path := filepath.Join(b.binariesFolder, binary.Name)

	file, err := os.Open(path)
	if err != nil {
		return server.BinaryData{}, fmt.Errorf("get: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return server.BinaryData{}, fmt.Errorf("get: %w", err)
	}

	return server.BinaryData{
		User:       binary.User,
		Name:       binary.Name,
		Metadata:   metadata,
		FileName:   binary.Filename,
		DataReader: file,
		Size:       stat.Size(),
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
			FileName: binary.Filename,
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
