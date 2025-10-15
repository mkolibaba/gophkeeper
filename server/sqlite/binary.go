package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
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

func (s *BinaryService) Save(ctx context.Context, data server.ReadableBinaryData, user string) error {
	// TODO: тут лучше проверить, что такой записи и такого файла нет

	err := s.qs.SaveBinary(ctx, sqlc.SaveBinaryParams{
		Name:     data.Name,
		Filename: data.Filename,
		Size:     data.Size,
		Notes:    stringOrNull(data.Notes),
		User:     user,
	})

	if err != nil {
		return tryUnwrapSaveError(err)
	}

	dest, err := os.Create(filepath.Join(s.binariesFolder, fmt.Sprintf("%s__%s", user, data.Name)))
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

func (s *BinaryService) Get(ctx context.Context, name string, user string) (*server.ReadableBinaryData, error) {
	binary, err := s.qs.GetBinary(ctx, name, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, server.ErrDataNotFound
		}
		return nil, fmt.Errorf("get: %w", err)
	}

	path := filepath.Join(s.binariesFolder, fmt.Sprintf("%s__%s", user, binary.Name))

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &server.ReadableBinaryData{
		BinaryData: server.BinaryData{
			Name:     binary.Name,
			Notes:    stringOrEmpty(binary.Notes),
			Filename: binary.Filename,
			Size:     binary.Size,
		},
		DataReader: file,
	}, nil
}

func (s *BinaryService) GetAll(ctx context.Context, user string) ([]server.BinaryData, error) {
	binaries, err := s.qs.GetAllBinaries(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.BinaryData
	for _, binary := range binaries {
		result = append(result, server.BinaryData{
			Name:     binary.Name,
			Filename: binary.Filename,
			Size:     binary.Size,
			Notes:    stringOrEmpty(binary.Notes),
		})
	}

	return result, nil
}

func (s *BinaryService) Update(ctx context.Context, data server.BinaryDataUpdate, user string) error {
	// TODO: implement
	panic("implement me")
}

func (s *BinaryService) Remove(ctx context.Context, name string, user string) error {
	n, err := s.qs.RemoveBinary(ctx, name)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if n == 0 {
		return server.ErrDataNotFound
	}
	path := filepath.Join(s.binariesFolder, fmt.Sprintf("%s__%s", user, name))
	if err = os.Remove(path); err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	return nil
}
