package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/sqlite/converter"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"io"
	"os"
	"path/filepath"
)

type BinaryService struct {
	qs             *sqlc.Queries
	converter      converter.DataConverter
	binariesFolder string
}

func NewBinaryService(queries *sqlc.Queries, db *DB, converter converter.DataConverter) *BinaryService {
	return &BinaryService{
		qs:             queries,
		converter:      converter,
		binariesFolder: db.binariesFolder,
	}
}

func (s *BinaryService) Create(ctx context.Context, data server.ReadableBinaryData) error {
	id, err := s.qs.InsertBinary(ctx, s.converter.ConvertToInsertBinary(ctx, data))
	if err != nil {
		return unwrapInsertError(err)
	}

	asset, err := os.Create(s.getBinaryAssetPath(id))
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}
	defer asset.Close()

	if _, err = io.CopyBuffer(asset, data.DataReader, make([]byte, 1024*1024)); err != nil {
		os.Remove(asset.Name())
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (s *BinaryService) Get(ctx context.Context, id int64) (*server.ReadableBinaryData, error) {
	binary, err := s.qs.SelectBinary(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, server.ErrDataNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	file, err := os.Open(s.getBinaryAssetPath(id))
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &server.ReadableBinaryData{
		BinaryData: s.converter.ConvertToBinaryData(binary),
		DataReader: file,
	}, nil
}

func (s *BinaryService) GetAll(ctx context.Context) ([]server.BinaryData, error) {
	return getAllData(ctx, s.qs.SelectBinaries, s.converter.ConvertToBinaryDataSlice)
}

func (s *BinaryService) Update(ctx context.Context, id int64, data server.BinaryDataUpdate) error {
	binary, err := s.qs.SelectBinary(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return server.ErrDataNotFound
	}
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	params := s.converter.ConvertToUpdateBinary(binary)
	s.converter.ConvertToUpdateBinaryUpdate(data, &params)

	n, err := s.qs.UpdateBinary(ctx, params)
	if n == 0 {
		return server.ErrDataNotFound
	}
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (s *BinaryService) Remove(ctx context.Context, id int64) error {
	if err := removeData(ctx, s.qs.DeleteBinary, id); err != nil {
		return err
	}

	if err := os.Remove(s.getBinaryAssetPath(id)); err != nil {
		return fmt.Errorf("remove: %w", err)
	}

	return nil
}

func (s *BinaryService) getBinaryAssetPath(id int64) string {
	return filepath.Join(s.binariesFolder, fmt.Sprintf("%d", id))
}
